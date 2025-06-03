package comment

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"
	usecase "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/comment"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase     *usecase.Usecase
	authService authpb.AuthServiceClient
}

func NewCommentHandler(r *gin.RouterGroup, uc *usecase.Usecase, authClient authpb.AuthServiceClient) {
	h := &Handler{
		usecase:     uc,
		authService: authClient,
	}

	r.GET("/comments", h.GetComments)
	r.POST("/comments/create", h.CreateComment)
	r.DELETE("/comments/delete", h.DeleteComment)
}

func (h *Handler) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}

	comments, err := h.usecase.GetCommentsByPost(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func (h *Handler) CreateComment(c *gin.Context) {
	var input struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token := extractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.authService.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{Token: token})
	log.Println("Проверка токена через gRPC в comment/handler.go:", token)
	if err != nil || !resp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	err = h.usecase.CreateComment(c.Request.Context(), input.PostID, resp.Username, input.Content)
	if err != nil {
		log.Println("Ошибка при создании комментария:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment created"})
}

func (h *Handler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Query("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment_id"})
		return
	}

	token := extractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.authService.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{Token: token})
	if err != nil || !resp.Valid || resp.Role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	err = h.usecase.DeleteComment(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

func extractBearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) > len(prefix) && header[:len(prefix)] == prefix {
		return header[len(prefix):]
	}
	return ""
}
