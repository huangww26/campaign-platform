package handler

import (
	"sync"
	"time"

	"campaign-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type RenderHandler struct {
	campaignSvc  *service.CampaignService
	componentSvc *service.ComponentService

	mu        sync.RWMutex
	validComps map[string]bool  // 组件名缓存
	lastRefresh time.Time
}

func NewRenderHandler(cs *service.CampaignService, compS *service.ComponentService) *RenderHandler {
	h := &RenderHandler{
		campaignSvc:  cs,
		componentSvc: compS,
		validComps:   make(map[string]bool),
	}
	h.refresh()
	return h
}

func (h *RenderHandler) refresh() {
	list, err := h.componentSvc.ListActive()
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.validComps = make(map[string]bool, len(list))
	for _, c := range list {
		h.validComps[c.Name] = true
	}
	h.lastRefresh = time.Now()
}

func (h *RenderHandler) isValidComponent(name string) bool {
	// 5 分钟刷新一次缓存
	if time.Since(h.lastRefresh) > 5*time.Minute {
		h.refresh()
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.validComps[name]
}

// --- 返回结构 ---

type RenderResp struct {
	Slug     string        `json:"slug"`
	Status   string        `json:"status"`
	Lang     string        `json:"lang"`
	Sections []SectionResp `json:"sections"`
}

type SectionResp struct {
	Component string                 `json:"component"`
	Props     map[string]interface{} `json:"props"`
}

// GET /api/v1/render/:slug
func (h *RenderHandler) RenderConfig(c *gin.Context) {
	slug := c.Param("slug")

	camp, err := h.campaignSvc.GetBySlug(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if camp.Status != "active" {
		c.JSON(404, gin.H{"error": "not available"})
		return
	}

	config := camp.Config

	sectionsRaw, ok := config["sections"]
	if !ok {
		c.JSON(200, RenderResp{Slug: camp.Slug, Status: camp.Status})
		return
	}

	sectionsList, ok := sectionsRaw.([]interface{})
	if !ok {
		c.JSON(500, gin.H{"error": "invalid sections"})
		return
	}

	var resp RenderResp
	resp.Slug = camp.Slug
	resp.Status = camp.Status
	if lang, ok := config["lang"].(string); ok {
		resp.Lang = lang
	}

	for _, s := range sectionsList {
		sec, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		compName, _ := sec["component"].(string)
		if compName == "" || !h.isValidComponent(compName) {
			continue
		}
		props, _ := sec["props"].(map[string]interface{})
		resp.Sections = append(resp.Sections, SectionResp{
			Component: compName,
			Props:     props,
		})
	}

	c.JSON(200, resp)
}
