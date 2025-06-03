package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/post"
	PostUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/post"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	uc *PostUC.UseCase
}

func NewPostHandler(r *gin.Engine, uc *PostUC.UseCase, authMiddleware gin.HandlerFunc) {
	h := &PostHandler{uc: uc}

	r.GET("/posts/all", h.getAll)
	r.GET("/posts", h.getByTopic)

	auth := r.Group("/", authMiddleware)
	auth.POST("/posts/create", h.create)
	auth.DELETE("/posts/delete", h.delete)
}

func (h *PostHandler) getAll(c *gin.Context) {
	posts, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Error getting all posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h *PostHandler) getByTopic(c *gin.Context) {
	topicIDStr := c.Query("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
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

func (h *PostHandler) create(c *gin.Context) {
	var req struct {
		TopicID int    `json:"topic_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&req); err != nil {
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

func (h *PostHandler) delete(c *gin.Context) {
	roleAny, _ := c.Get("role")
	if roleAny.(string) != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete posts"})
		return
	}

	postIDStr := c.Query("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}

	err = h.uc.Delete(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
