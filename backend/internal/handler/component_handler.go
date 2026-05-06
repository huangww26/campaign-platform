package handler

import (
	"campaign-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type ComponentHandler struct {
	svc *service.ComponentService
}

func NewComponentHandler(svc *service.ComponentService) *ComponentHandler {
	return &ComponentHandler{svc: svc}
}

// GET /api/v1/components
func (h *ComponentHandler) ListComponents(c *gin.Context) {
	list, err := h.svc.ListActive()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}
