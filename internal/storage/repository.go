package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"time"
)

type Authorization interface {
	SaveUser(user models.User) (int64, error)
	GetUser(email string) (models.User, error)
	IsAdmin(id int64) (bool, error)
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
	GetModuleData(moduleId int64, startDate time.Time, endDate time.Time) ([]models.ModuleData, error)
	SaveModuleData(moduleData models.ModuleData) error
}

type Repository struct {
	Authorization
	PlantFamily
	Module
	Settings
	Recommendations
	ModuleData
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization:   NewAuthPostgres(db),
		PlantFamily:     NewPlantFamilyPostgres(db),
		Module:          NewModulePostgres(db),
		Settings:        NewSettingsPostgres(db),
		Recommendations: NewRecommendationsPostgres(db),
		ModuleData:      NewModuleDataPostgres(db),
	}
}
