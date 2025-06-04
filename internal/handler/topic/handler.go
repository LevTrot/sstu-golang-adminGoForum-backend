package handler

import (
	"net/http"
	"strconv"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/topic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TopicHandler struct {
	UseCase *topic.UseCase
	logger  *zap.Logger
}

func NewTopicHandler(rg *gin.RouterGroup, uc *topic.UseCase, authMiddleware gin.HandlerFunc, logger *zap.Logger) {
	h := &TopicHandler{UseCase: uc, logger: logger}

	rg.GET("/topics", h.GetAll)
	rg.POST("/topics/create", authMiddleware, h.RequireAdmin(), h.Create)
	rg.DELETE("/topics/delete", authMiddleware, h.RequireAdmin(), h.Delete)
}

func (h *TopicHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	rg.GET("/topics", h.GetAll)
	rg.POST("/topics/create", authMiddleware, h.RequireAdmin(), h.Create)
	rg.DELETE("/topics/delete", authMiddleware, h.RequireAdmin(), h.Delete)
}

// GetAll godoc
// @Summary Get all topics
// @Tags Topics
// @Produce json
// @Success 200 {array} topic.Topic
// @Failure 500 {object} response.ErrorResponse
// @Router /topics [get]
func (h *TopicHandler) GetAll(c *gin.Context) {
	topics, err := h.UseCase.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get topics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении тем"})
		return
	}
	c.JSON(http.StatusOK, topics)
}

type CreateTopicInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// Create godoc
// @Summary Create a new topic
// @Tags Topics
// @Accept json
// @Produce json
// @Param topic body CreateTopicInput true "Topic input"
// @Success 201 {object} response.MessageResponse
// @Failure 400,500 {object} response.ErrorResponse
// @Router /topics/create [post]
func (h *TopicHandler) Create(c *gin.Context) {
	var input CreateTopicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid topic input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	err := h.UseCase.Create(c.Request.Context(), input.Title, input.Description)
	if err != nil {
		h.logger.Error("invalid topic input", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании темы"})
		return
	}
	c.Status(http.StatusCreated)
}

// Delete godoc
// @Summary Delete a topic by ID (admin only)
// @Tags Topics
// @Produce json
// @Param id query int true "Topic ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400,500 {object} response.ErrorResponse
// @Router /topics/delete [delete]
func (h *TopicHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid topic ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	err = h.UseCase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении темы"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TopicHandler) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRaw, exists := c.Get("role")
		if !exists || roleRaw != "ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён: требуется роль ADMIN"})
			return
		}
		c.Next()
	}
}
