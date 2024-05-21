package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
)

type Authorization interface {
	SaveUser(user models.User) (int64, error)
	GetUser(email string) (models.User, error)
}

type PlantFamily interface {
}

type Module interface {
}

type Repository struct {
	Authorization
	PlantFamily
	Module
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
