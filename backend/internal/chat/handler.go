package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles chat HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new chat handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers chat routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.Chat)
}

// Chat handles chat requests
func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Chat(c.Request.Context(), req.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
