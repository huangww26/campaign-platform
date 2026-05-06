package service

import (
	"errors"
	"fmt"

	"campaign-platform/internal/model"
	"campaign-platform/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidStatus = errors.New("invalid status transition")
	ErrInvalidConfig = errors.New("invalid config")
)

// 合法状态流转
var validTransitions = map[string][]string{
	"draft":  {"active"},
	"active": {"ended"},
	"ended":  {},
}

type CampaignService struct {
	repo *repository.CampaignRepo
}

func NewCampaignService(repo *repository.CampaignRepo) *CampaignService {
	return &CampaignService{repo: repo}
}

// --- 模板 ---

func (s *CampaignService) ListTemplates() ([]model.CampaignTemplate, error) {
	return s.repo.ListTemplates()
}

func (s *CampaignService) CreateTemplate(t *model.CampaignTemplate) error {
	return s.repo.CreateTemplate(t)
}

// --- 活动 ---

func (s *CampaignService) ListCampaigns(status string) ([]model.Campaign, error) {
	return s.repo.ListCampaigns(status)
}

func (s *CampaignService) GetByID(id uint) (*model.Campaign, error) {
	c, err := s.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return c, err
}

func (s *CampaignService) GetBySlug(slug string) (*model.Campaign, error) {
	c, err := s.repo.GetBySlug(slug)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return c, err
}

func (s *CampaignService) Create(c *model.Campaign) error {
	c.Version = 1
	c.Status = "draft"
	return s.repo.Create(c)
}

func (s *CampaignService) Update(id uint, config model.JSON, changelog string) (*model.Campaign, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 保存当前版本到历史
	old := model.CampaignVersion{
		CampaignID: c.ID,
		Version:    c.Version,
		Config:     c.Config,
		Changelog:  changelog,
	}
	if err := s.repo.CreateVersion(&old); err != nil {
		return nil, fmt.Errorf("save version: %w", err)
	}

	// 更新配置，版本号 +1
	c.Config = config
	c.Version++
	if err := s.repo.Update(c); err != nil {
		return nil, fmt.Errorf("update campaign: %w", err)
	}

	return c, nil
}

func (s *CampaignService) UpdateStatus(id uint, newStatus string) (*model.Campaign, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	allowed, ok := validTransitions[c.Status]
	if !ok {
		return nil, ErrInvalidStatus
	}
	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}
	if !valid {
		return nil, ErrInvalidStatus
	}

	c.Status = newStatus
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CampaignService) ListVersions(campaignID uint) ([]model.CampaignVersion, error) {
	return s.repo.ListVersions(campaignID)
}
