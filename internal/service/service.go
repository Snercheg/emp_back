package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
)

type Authorization interface {
	SaveUser(user models.User) (int64, error)
	GenerateToken(email, password string) (string, error)
}

type PlantFamily interface {
}

type Module interface {
}

type Service struct {
	Authorization
	PlantFamily
	Module
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
