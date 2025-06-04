package comment

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"
	usecase "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/comment"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase     *usecase.Usecase
	authService authpb.AuthServiceClient
	logger      *zap.Logger
}

type CreateCommentInput struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

func NewCommentHandler(r *gin.RouterGroup, uc *usecase.Usecase, authClient authpb.AuthServiceClient, logger *zap.Logger) {
	h := &Handler{
		usecase:     uc,
		authService: authClient,
		logger:      logger,
	}

	r.GET("/comments", h.GetComments)
	r.POST("/comments/create", h.CreateComment)
	r.DELETE("/comments/delete", h.DeleteComment)
}

// GetComments godoc
// @Summary Get comments for a post
// @Tags Comments
// @Produce json
// @Param post_id query int true "Post ID"
// @Success 200 {object} response.DataCommentsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /comments [get]
func (h *Handler) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil {
		h.logger.Error("invalid post_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}

	comments, err := h.usecase.GetCommentsByPost(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("failed to fetch comments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// CreateComment godoc
// @Summary Create a new comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body comment.CreateCommentInput true "Comment content"
// @Success 200 {object} response.MessageResponse
// @Failure 400,401,500 {object} response.ErrorResponse
// @Router /comments/create [post]
func (h *Handler) CreateComment(c *gin.Context) {
	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token := extractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.authService.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{Token: token})
	if err != nil || !resp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	err = h.usecase.CreateComment(c.Request.Context(), input.PostID, resp.Username, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment created"})
}

// DeleteComment godoc
// @Summary Delete a comment by ID (admin only)
// @Tags Comments
// @Produce json
// @Param comment_id query int true "Comment ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400,401,403,500 {object} response.ErrorResponse
// @Router /comments/delete [delete]
func (h *Handler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Query("comment_id"))
	if err != nil {
		h.logger.Error("invalid comment_id", zap.Error(err))
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
		h.logger.Error("invalid delete", zap.Error(err))
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
