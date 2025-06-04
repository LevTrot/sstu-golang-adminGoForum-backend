package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"
	domain "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/chat"
	chatUsecase "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ChatHandler struct {
	usecase     *chatUsecase.UseCase
	authService authpb.AuthServiceClient
	clients     map[*websocket.Conn]bool
	broadcast   chan domain.ChatMessage
	logger      *zap.Logger
}

func New(usecase *chatUsecase.UseCase, authService authpb.AuthServiceClient, logger *zap.Logger) *ChatHandler {
	h := &ChatHandler{
		usecase:     usecase,
		authService: authService,
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan domain.ChatMessage),
		logger:      logger,
	}
	go h.handleMessages()
	return h
}

// GetMessagesHandler godoc
// @Summary Get recent chat messages
// @Tags Chat
// @Produce json
// @Success 200 {array} domain.ChatMessage
// @Failure 500 {object} response.ErrorResponse
// @Router /chat/messages [get]
func (h *ChatHandler) GetMessagesHandler(c *gin.Context) {
	messages, err := h.usecase.GetMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ChatWebSocketHandler godoc
// @Summary WebSocket endpoint for real-time chat
// @Tags Chat
// @Produce plain
// @Param token query string true "JWT token"
// @Success 101 {string} string "WebSocket Connection Established"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /chat [get]
func (h *ChatHandler) ChatWebSocketHandler(c *gin.Context) {

	token := c.Query("token")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := h.authService.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		h.logger.Error("Failed:", zap.Error(err))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !resp.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	username := resp.Username

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.clients[conn] = true

	go func(conn *websocket.Conn, username string) {
		defer func() {
			conn.Close()
			delete(h.clients, conn)
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var payload struct {
				Content string `json:"content"`
			}
			if err := json.Unmarshal(msg, &payload); err != nil {
				continue
			}
			if strings.TrimSpace(payload.Content) == "" {
				continue
			}

			message := domain.ChatMessage{
				Username:  username,
				Content:   payload.Content,
				Timestamp: time.Now(),
			}
			if err := h.usecase.SendMessage(context.Background(), message.Username, message.Content); err != nil {
				continue
			}

			h.broadcast <- message
		}
	}(conn, username)
}

func (h *ChatHandler) handleMessages() {
	for {
		msg := <-h.broadcast
		data, _ := json.Marshal(msg)
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
