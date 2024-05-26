package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type RecommendationsPostgres struct {
	DB *sql.DB
}

func NewRecommendationsPostgres(db *sql.DB) *RecommendationsPostgres {
	return &RecommendationsPostgres{DB: db}
}

func (r *RecommendationsPostgres) GetRecommendation(id int64) (models.Recommendation, error) {
	op := "storage.repository.GetRecommendation"
	row := r.DB.QueryRow("SELECT * FROM recommendations WHERE id = $1", id)

	var recommendation models.Recommendation
	err := row.Scan(&recommendation.ID, &recommendation.Title, &recommendation.TemperatureMin, &recommendation.TemperatureMax, &recommendation.HumidityInMin, &recommendation.HumidityInMax, &recommendation.HumidityOutMin, &recommendation.HumidityOutMax, &recommendation.IlluminationMin, &recommendation.IlluminationMax, &recommendation.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recommendation, ErrRecommendationNotFound
		}
		return recommendation, fmt.Errorf("%s: %v", op, err)
	}
	return recommendation, nil
}

func (r *RecommendationsPostgres) GetRecommendations() ([]models.Recommendation, error) {
	op := "storage.repository.GetRecommendations"
	rows, err := r.DB.Query("SELECT * FROM recommendations")
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()
	recommendations := []models.Recommendation{}
	for rows.Next() {
		var recommendation models.Recommendation
		err := rows.Scan(&recommendation.ID, &recommendation.Title, &recommendation.TemperatureMin, &recommendation.TemperatureMax, &recommendation.HumidityInMin, &recommendation.HumidityInMax, &recommendation.HumidityOutMin, &recommendation.HumidityOutMax, &recommendation.IlluminationMin, &recommendation.IlluminationMax, &recommendation.Description)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrRecommendationNotFound
			}
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		recommendations = append(recommendations, recommendation)
	}
	return recommendations, nil
}

func (r *RecommendationsPostgres) SaveRecommendation(recommendation models.Recommendation) (int64, error) {
	op := "storage.repository.SaveRecommendation"
	result, err := r.DB.Exec("INSERT INTO recommendations (title, temperature_min, temperature_max, humidity_in_min, humidity_in_max, humidity_out_min, humidity_out_max, illumination_min, illumination_max, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", recommendation.Title, recommendation.TemperatureMin, recommendation.TemperatureMax, recommendation.HumidityInMin, recommendation.HumidityInMax, recommendation.HumidityOutMin, recommendation.HumidityOutMax, recommendation.IlluminationMin, recommendation.IlluminationMax, recommendation.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, ErrRecommendationExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

func (r *RecommendationsPostgres) UpdateRecommendation(id int64, recommendation models.Recommendation) error {
	op := "storage.repository.UpdateRecommendation"
	result, err := r.DB.Exec("UPDATE recommendations SET title = $1, temperature_min = $2, temperature_max = $3, humidity_in_min = $4, humidity_in_max = $5, humidity_out_min = $6, humidity_out_max = $7, illumination_min = $8, illumination_max = $9, description = $10 WHERE id = $11", recommendation.Title, recommendation.TemperatureMin, recommendation.TemperatureMax, recommendation.HumidityInMin, recommendation.HumidityInMax, recommendation.HumidityOutMin, recommendation.HumidityOutMax, recommendation.IlluminationMin, recommendation.IlluminationMax, recommendation.Description, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, ErrRecommendationCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	if rowsAffected == 0 {
		return ErrRecommendationNotFound
	}
	return nil
}

func (r *RecommendationsPostgres) DeleteRecommendation(id int64) error {
	op := "storage.repository.DeleteRecommendation"
	result, err := r.DB.Exec("DELETE FROM recommendations WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key
			return fmt.Errorf("%s: %v", op, err, ErrRecommendationNotFound)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	if rowsAffected == 0 {
		return ErrRecommendationNotFound
	}
	return nil
}
