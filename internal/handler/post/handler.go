package handler

import (
	"net/http"
	"strconv"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/post"
	PostUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/post"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	uc     *PostUC.UseCase
	logger *zap.Logger
}

type CreatePostInput struct {
	TopicID int    `json:"topic_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewPostHandler(r *gin.Engine, uc *PostUC.UseCase, authMiddleware gin.HandlerFunc, logger *zap.Logger) {
	h := &PostHandler{uc: uc, logger: logger}

	r.GET("/posts/all", h.getAll)
	r.GET("/posts", h.getByTopic)

	auth := r.Group("/", authMiddleware)
	auth.POST("/posts/create", h.create)
	auth.DELETE("/posts/delete", h.delete)
}

// getAll godoc
// @Summary Get all posts
// @Tags Posts
// @Produce json
// @Success 200 {object} response.DataPostsResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posts/all [get]
func (h *PostHandler) getAll(c *gin.Context) {
	posts, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// getByTopic godoc
// @Summary Get posts by topic
// @Tags Posts
// @Produce json
// @Param topic_id query int true "Topic ID"
// @Success 200 {object} response.DataPostsResponse
// @Failure 400,500 {object} response.ErrorResponse
// @Router /posts [get]
func (h *PostHandler) getByTopic(c *gin.Context) {
	topicIDStr := c.Query("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		h.logger.Error("invalid topic_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic_id"})
		return
	}

	posts, err := h.uc.GetByTopic(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// create godoc
// @Summary Create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body CreatePostInput true "Post payload"
// @Success 200 {object} response.MessageResponse
// @Failure 400,500 {object} response.ErrorResponse
// @Router /posts/create [post]
func (h *PostHandler) create(c *gin.Context) {
	var req CreatePostInput
	if err := c.BindJSON(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	usernameAny, _ := c.Get("username")
	username := usernameAny.(string)

	p := post.Post{
		TopicID:  req.TopicID,
		Title:    req.Title,
		Content:  req.Content,
		Username: username,
	}

	err := h.uc.Create(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post created"})
}

// delete godoc
// @Summary Delete post by ID (admin only)
// @Tags Posts
// @Produce json
// @Param post_id query int true "Post ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400,403,500 {object} response.ErrorResponse
// @Router /posts/delete [delete]
func (h *PostHandler) delete(c *gin.Context) {
	roleAny, _ := c.Get("role")
	if roleAny.(string) != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete posts"})
		return
	}

	postIDStr := c.Query("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Error("invalid post_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}

	err = h.uc.Delete(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("failed to delete post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
