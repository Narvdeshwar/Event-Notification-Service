package api

import "github.com/gin-gonic/gin"
import "github.com/prometheus/client_golang/prometheus/promhttp"



func (h *Handler) CreateEvent(c *gin.Context) {
	h.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
