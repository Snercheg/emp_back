package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
)

type ModuleService struct {
	repo storage.Module
	log  *slog.Logger
}

func NewModuleService(repo storage.Module, log *slog.Logger) *ModuleService {
	return &ModuleService{repo: repo, log: log}
}

func (s *ModuleService) SaveModule(userId int64, module models.Module) (int64, error) {
	return s.repo.SaveModule(userId, module)
}

func (s *ModuleService) GetModules(userID int64) ([]models.Module, error) {
	return s.repo.GetModules(userID)
}

func (s *ModuleService) GetModule(userId, moduleId int64) (*models.Module, error) {
	return s.repo.GetModule(userId, moduleId)
}

func (s *ModuleService) DeleteModule(userId, moduleId int64) error {
	return s.repo.DeleteModule(userId, moduleId)
}

func (s *ModuleService) UpdateModule(userId, moduleId int64, module models.UpdateModuleInput) error {
	if err := module.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateModule(userId, moduleId, module)
}
