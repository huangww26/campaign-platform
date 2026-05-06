package repository

import (
	"campaign-platform/internal/model"

	"gorm.io/gorm"
)

type CampaignRepo struct {
	db *gorm.DB
}

func NewCampaignRepo(db *gorm.DB) *CampaignRepo {
	return &CampaignRepo{db: db}
}

// --- 模板 ---

func (r *CampaignRepo) ListTemplates() ([]model.CampaignTemplate, error) {
	var list []model.CampaignTemplate
	err := r.db.Order("id").Find(&list).Error
	return list, err
}

func (r *CampaignRepo) CreateTemplate(t *model.CampaignTemplate) error {
	return r.db.Create(t).Error
}

// --- 活动 ---

func (r *CampaignRepo) ListCampaigns(status string) ([]model.Campaign, error) {
	var list []model.Campaign
	q := r.db.Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *CampaignRepo) GetByID(id uint) (*model.Campaign, error) {
	var c model.Campaign
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *CampaignRepo) GetBySlug(slug string) (*model.Campaign, error) {
	var c model.Campaign
	err := r.db.Where("slug = ?", slug).First(&c).Error
	return &c, err
}

func (r *CampaignRepo) Create(c *model.Campaign) error {
	return r.db.Create(c).Error
}

func (r *CampaignRepo) Update(c *model.Campaign) error {
	return r.db.Save(c).Error
}

// --- 版本 ---

func (r *CampaignRepo) ListVersions(campaignID uint) ([]model.CampaignVersion, error) {
	var list []model.CampaignVersion
	err := r.db.Where("campaign_id = ?", campaignID).
		Order("version DESC").Find(&list).Error
	return list, err
}

func (r *CampaignRepo) CreateVersion(v *model.CampaignVersion) error {
	return r.db.Create(v).Error
}
