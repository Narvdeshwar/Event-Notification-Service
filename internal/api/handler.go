package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"notification-service/internal/service"
)

type Handler struct {
	service *service.NotificationService
}

func New(svc *service.NotificationService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req struct {
		Type      string `json:"type" binding:"required"`
		Recipient string `json:"recipient" binding:"required"`
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}
