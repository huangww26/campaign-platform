package handler

import (
	"net/http"
	"strconv"

	"campaign-platform/internal/model"
	"campaign-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	svc *service.CampaignService
}

func NewCampaignHandler(svc *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{svc: svc}
}

// GET /api/v1/templates
func (h *CampaignHandler) ListTemplates(c *gin.Context) {
	list, err := h.svc.ListTemplates()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// POST /api/v1/templates
func (h *CampaignHandler) CreateTemplate(c *gin.Context) {
	var t model.CampaignTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateTemplate(&t); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, t)
}

// GET /api/v1/campaigns
func (h *CampaignHandler) ListCampaigns(c *gin.Context) {
	status := c.Query("status")
	list, err := h.svc.ListCampaigns(status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// GET /api/v1/campaigns/:id
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	camp, err := h.svc.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, camp)
}

// POST /api/v1/campaigns
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var camp model.Campaign
	if err := c.ShouldBindJSON(&camp); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(&camp); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, camp)
}

// PUT /api/v1/campaigns/:id
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Config    model.JSON `json:"config"`
		Changelog string     `json:"changelog"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	camp, err := h.svc.Update(uint(id), body.Config, body.Changelog)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, camp)
}

// PATCH /api/v1/campaigns/:id/status
func (h *CampaignHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	camp, err := h.svc.UpdateStatus(uint(id), body.Status)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		if err == service.ErrInvalidStatus {
			c.JSON(400, gin.H{"error": "invalid status transition"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, camp)
}

// GET /api/v1/campaigns/:id/versions
func (h *CampaignHandler) ListVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	list, err := h.svc.ListVersions(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
