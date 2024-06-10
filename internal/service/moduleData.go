package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
	"time"
)

type ModuleDataService struct {
	repo storage.ModuleData
	log  *slog.Logger
}

func NewModuleDataService(repo storage.ModuleData, log *slog.Logger) *ModuleDataService {
	return &ModuleDataService{repo: repo, log: log}
}

func (s *ModuleDataService) SaveModuleData(moduleData models.ModuleData) error {
	return s.repo.SaveModuleData(moduleData)
}

func (s *ModuleDataService) GetModuleData(moduleId int64, startDate, endDate time.Time) ([]models.ModuleData, error) {
	return s.repo.GetModuleData(moduleId, startDate, endDate)
}
