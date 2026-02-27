package api

import (
	"event-driven-notification-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.NotificationService
}

func New(s *service.NotificationService) *Handler {
	return &Handler{
		service: s,
	}
}
func (h *Handler) CreateEvent(c *gin.Context) {

	var req struct {
		Type      string `json:"type"`
		Recipient string `json:"recipient"`
		Payload   string `json:"payload"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Enqueue(
		c.Request.Context(),
		req.Type,
		req.Recipient,
		[]byte(req.Payload),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}
