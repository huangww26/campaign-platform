package service

import (
	"campaign-platform/internal/model"
	"campaign-platform/internal/repository"
)

type ComponentService struct {
	repo *repository.ComponentRepo
}

func NewComponentService(repo *repository.ComponentRepo) *ComponentService {
	return &ComponentService{repo: repo}
}

func (s *ComponentService) ListActive() ([]model.Component, error) {
	return s.repo.ListActive()
}
