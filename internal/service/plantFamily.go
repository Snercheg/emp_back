package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
)

type PlantFamilyService struct {
	repo storage.PlantFamily
	log  *slog.Logger
}

func NewPlantFamilyService(repo storage.PlantFamily, log *slog.Logger) *PlantFamilyService {
	return &PlantFamilyService{repo: repo, log: log}
}

func (s *PlantFamilyService) SavePlantFamily(plantFamily models.PlantFamily) (int64, error) {
	return s.repo.SavePlantFamily(plantFamily)
}

func (s *PlantFamilyService) GetPlantFamily(id int64) (*models.PlantFamily, error) {
	return s.repo.GetPlantFamily(id)
}

func (s *PlantFamilyService) GetPlantFamilies() ([]models.PlantFamily, error) {
	return s.repo.GetPlantFamilies()
}

func (s *PlantFamilyService) DeletePlantFamily(id int64) error {
	return s.repo.DeletePlantFamily(id)
}

func (s *PlantFamilyService) UpdatePlantFamily(id int64, plantFamily models.PlantFamily) error {
	return s.repo.UpdatePlantFamily(id, plantFamily)
}
