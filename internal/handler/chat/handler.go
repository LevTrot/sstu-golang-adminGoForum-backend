package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"
	domain "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/chat"
	chatUsecase "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	usecase     *chatUsecase.UseCase
	authService authpb.AuthServiceClient
	clients     map[*websocket.Conn]bool
	broadcast   chan domain.ChatMessage
}

func New(usecase *chatUsecase.UseCase, authService authpb.AuthServiceClient) *ChatHandler {
	h := &ChatHandler{
		usecase:     usecase,
		authService: authService,
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan domain.ChatMessage),
	}
	go h.handleMessages()
	return h
}

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

func (h *ChatHandler) ChatWebSocketHandler(c *gin.Context) {
	log.Println("Получен новый WS-запрос на /chat")

	token := c.Query("token")
	if token == "" {
		log.Println("Нет токена в query")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := h.authService.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		log.Println("Ошибка при валидации токена:", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !resp.Valid {
		log.Println("Невалидный токен:", resp.GetError())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	username := resp.Username
	log.Println("Успешная валидация токена. Пользователь:", username)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ошибка апгрейда WebSocket:", err)
		return
	}
	log.Println("WebSocket соединение установлено для:", username)
	h.clients[conn] = true

	go func(conn *websocket.Conn, username string) {
		defer func() {
			conn.Close()
			delete(h.clients, conn)
			log.Println("Соединение закрыто для:", username)
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения сообщения от", username, ":", err)
				break
			}

			var payload struct {
				Content string `json:"content"`
			}
			if err := json.Unmarshal(msg, &payload); err != nil {
				log.Println("Невалидный JSON от", username, ":", err)
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
				log.Println("Ошибка сохранения сообщения:", err)
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
				log.Println("Ошибка отправки сообщения клиенту:", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
