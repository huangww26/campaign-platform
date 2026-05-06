package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSON 自定义类型，支持 GORM JSONB
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// ==== 数据表 ====

// CampaignTemplate 模板定义
type CampaignTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	SchemaDef JSON      `gorm:"type:jsonb;default:'{}'" json:"schema_def"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Campaign 活动实例
type Campaign struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TemplateID *uint     `json:"template_id,omitempty"`
	Name       string    `gorm:"size:128;not null" json:"name"`
	Slug       string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Status     string    `gorm:"size:16;default:draft" json:"status"`
	Config     JSON      `gorm:"type:jsonb;default:'{}'" json:"config"`
	Version    int       `gorm:"default:1" json:"version"`
	CreatedBy  string    `gorm:"size:64" json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Component 组件注册表
type Component struct {
	Name        string `gorm:"primaryKey;size:64" json:"name"`
	Version     string `gorm:"size:16;default:v1" json:"version"`
	Category    string `gorm:"size:16" json:"category"`
	PropsSchema JSON   `gorm:"type:jsonb;default:'{}'" json:"props_schema"`
	Status      string `gorm:"size:16;default:active" json:"status"`
}

// CampaignVersion 版本历史
type CampaignVersion struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CampaignID uint       `gorm:"index;not null" json:"campaign_id"`
	Version    int        `gorm:"not null" json:"version"`
	Config     JSON       `gorm:"type:jsonb;not null" json:"config"`
	Changelog  string     `gorm:"type:text" json:"changelog,omitempty"`
	DeployedAt *time.Time `json:"deployed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
