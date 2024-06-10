package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
)

type SettingsService struct {
	repo storage.Settings
	log  *slog.Logger
}

func NewSettingsService(repo storage.Settings, log *slog.Logger) *SettingsService {
	return &SettingsService{repo: repo, log: log}
}

func (s *SettingsService) SaveSetting(moduleId int64, settings models.Setting) (int64, error) {
	return s.repo.SaveSetting(moduleId, settings)
}

func (s *SettingsService) GetSetting(moduleId int64) (*models.Setting, error) {
	return s.repo.GetSetting(moduleId)
}

func (s *SettingsService) UpdateSetting(moduleId int64, setting models.Setting) error {
	return s.repo.UpdateSetting(moduleId, setting)
}
