package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
	"time"
)

type Authorization interface {
	SaveUser(user models.User) (int64, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int64, error)
	IsAdmin(userId int64) (bool, error)
}

type PlantFamily interface {
	SavePlantFamily(plantFamily models.PlantFamily) (int64, error)
	GetPlantFamily(id int64) (*models.PlantFamily, error)
	GetPlantFamilies() ([]models.PlantFamily, error)
	DeletePlantFamily(id int64) error
	UpdatePlantFamily(id int64, plantFamily models.PlantFamily) error
}

type Module interface {
	SaveModule(userId int64, module models.Module) (int64, error)
	GetModules(userId int64) ([]models.Module, error)
	GetModule(userId, moduleId int64) (*models.Module, error)
	DeleteModule(userId, moduleId int64) error
	UpdateModule(userId, moduleId int64, module models.UpdateModuleInput) error
}
type Settings interface {
	SaveSetting(moduleId int64, settings models.Setting) (int64, error)
	GetSetting(moduleId int64) (*models.Setting, error)
	UpdateSetting(moduleId int64, setting models.Setting) error
}
type Recommendations interface {
	GetRecommendation(id int64) (models.Recommendation, error)
	GetRecommendations() ([]models.Recommendation, error)
	SaveRecommendation(recommendation models.Recommendation) (int64, error)
	UpdateRecommendation(id int64, recommendation models.Recommendation) error
	DeleteRecommendation(id int64) error
}
type ModuleData interface {
	GetModuleData(moduleId int64, startDate, endDate time.Time) ([]models.ModuleData, error)
	SaveModuleData(moduleData models.ModuleData) error
}

type Service struct {
	Authorization
	PlantFamily
	Module
	Settings
	Recommendations
	ModuleData
}

func NewService(repos *storage.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization:   NewAuthService(repos.Authorization, log),
		PlantFamily:     NewPlantFamilyService(repos.PlantFamily, log),
		Module:          NewModuleService(repos.Module, log),
		Settings:        NewSettingsService(repos.Settings, log),
		Recommendations: NewRecommendationsService(repos.Recommendations, log),
		ModuleData:      NewModuleDataService(repos.ModuleData, log),
	}
}
