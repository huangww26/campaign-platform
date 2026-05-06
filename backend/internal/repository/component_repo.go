package repository

import (
	"campaign-platform/internal/model"

	"gorm.io/gorm"
)

type ComponentRepo struct {
	db *gorm.DB
}

func NewComponentRepo(db *gorm.DB) *ComponentRepo {
	return &ComponentRepo{db: db}
}

func (r *ComponentRepo) ListActive() ([]model.Component, error) {
	var list []model.Component
	err := r.db.Where("status = ?", "active").Order("name").Find(&list).Error
	return list, err
}
