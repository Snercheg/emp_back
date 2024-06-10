package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type PlantFamilyPostgres struct {
	DB *sql.DB
}

func NewPlantFamilyPostgres(db *sql.DB) *PlantFamilyPostgres {
	return &PlantFamilyPostgres{DB: db}
}

func (p *PlantFamilyPostgres) SavePlantFamily(plantFamily models.PlantFamily) (int64, error) {
	op := "storage.repository.SaveUser"
	stmt, err := p.DB.Prepare(fmt.Sprintf("INSERT INTO plant_family (name, recommendation_id, picture_url, description) VALUES ($1, $2, $3, $4) RETURNING id", plantFamily.Name, plantFamily.RecommendationId, plantFamily.PictureUrl, plantFamily.Description))
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	res, err := stmt.Exec(plantFamily.Name, plantFamily.RecommendationId, plantFamily.PictureUrl, plantFamily.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, ErrPlantFamilyExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

func (p *PlantFamilyPostgres) GetPlantFamily(id int64) (*models.PlantFamily, error) {
	op := "storage.repository.GetPlantFamily"
	stmt, err := p.DB.Prepare(fmt.Sprintf("SELECT id, name, recommendation_id, picture_url, description FROM plant_family WHERE id = $1", id))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	row := stmt.QueryRow(id)
	plantFamily := &models.PlantFamily{}
	err = row.Scan(&plantFamily.ID, &plantFamily.Name, &plantFamily.RecommendationId, &plantFamily.PictureUrl, &plantFamily.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return plantFamily, ErrPlantFamilyNotFound
		}
		return plantFamily, fmt.Errorf("%s: %v", op, err)
	}
	return plantFamily, nil
}

func (p *PlantFamilyPostgres) GetPlantFamilies() ([]models.PlantFamily, error) {
	op := "storage.repository.GetPlantFamilies"
	stmt, err := p.DB.Prepare("SELECT id, name, recommendation_id, picture_url, description FROM plant_family")
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	var plantFamilies []models.PlantFamily
	for rows.Next() {
		plantFamily := models.PlantFamily{}
		err = rows.Scan(&plantFamily.ID, &plantFamily.Name, &plantFamily.RecommendationId, &plantFamily.PictureUrl, &plantFamily.Description)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		plantFamilies = append(plantFamilies, plantFamily)
	}
	return plantFamilies, nil
}

func (p *PlantFamilyPostgres) DeletePlantFamily(id int64) error {
	op := "storage.repository.DeletePlantFamily"
	stmt, err := p.DB.Prepare(fmt.Sprintf("DELETE FROM plant_family WHERE id = $1", id))
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key
			return fmt.Errorf("%s: %v", op, err, ErrPlantFamilyNotFound)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (p *PlantFamilyPostgres) UpdatePlantFamily(id int64, plantFamily models.PlantFamily) error {
	op := "storage.repository.UpdatePlantFamily"
	stmt, err := p.DB.Prepare(fmt.Sprintf("UPDATE plant_family SET name = $1, recommendation_id = $2, picture_url = $3, description = $4 WHERE id = $5", plantFamily.Name, plantFamily.RecommendationId, plantFamily.PictureUrl, plantFamily.Description, id))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, ErrPlantFamilyCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.Exec(plantFamily.Name, plantFamily.RecommendationId, plantFamily.PictureUrl, plantFamily.Description)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
