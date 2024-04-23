package postgres

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"

	"os"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	op := "storage.postgres.New"

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username string, email string, passHash []byte) (int64, error) {
	op := "storage.postgres.SaveUser"

	_, err := s.db.Prepare("INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := s.db.ExecContext(ctx, username, email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrUserExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	// Return the ID of the inserted row.

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// User returns a user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	op := "storage.postgres.User"
	row := s.db.QueryRowContext(ctx, "SELECT id, username, email, pass_hash FROM users WHERE email = $1", email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, storage.ErrUserNotFound
		}
		return u, fmt.Errorf("%s: %v", op, err)
	}
	return u, nil
}

// IsAdmin returns true if the user is an admin.
func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	op := "storage.postgres.IsAdmin"
	row := s.db.QueryRowContext(ctx, "SELECT is_admin FROM users WHERE id = $1", id)
	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return isAdmin, storage.ErrAppNotFound
		}
		return isAdmin, fmt.Errorf("%s: %v", op, err)
	}
	return isAdmin, nil
}

// SaveRecommendation saves a new plant recommendation.
func (s *Storage) SaveRecommendation(ctx context.Context, r *models.Recommendation) (int64, error) {
	op := "storage.postgres.SaveRecommendation"
	_, err := s.db.Prepare("INSERT INTO recommendations (title, temperature_min, temperature_max, humidity_min, humidity_max, illumination_min, illumination_max, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	res, err := s.db.ExecContext(ctx, r.Title, r.TemperatureMin, r.TemperatureMax, r.HumidityMin, r.HumidityMax, r.IlluminationMin, r.IlluminationMax, r.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrRecommendationExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	// Return the ID of the inserted row.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// Recommendation returns a recommendation by ID.
func (s *Storage) Recommendation(ctx context.Context, id int64) (models.Recommendation, error) {
	op := "storage.postgres.Recommendation"
	row := s.db.QueryRowContext(ctx, "SELECT id, title, temperature_min, temperature_max, humidity_min, humidity_max, illumination_min, illumination_max, description FROM recommendations WHERE id = $1", id)
	var r models.Recommendation
	err := row.Scan(&r.ID, &r.Title, &r.TemperatureMin, &r.TemperatureMax, &r.HumidityMin, &r.HumidityMax, &r.IlluminationMin, &r.IlluminationMax, &r.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, storage.ErrRecommendationNotFound
		}
		return r, fmt.Errorf("%s: %v", op, err)
	}
	return r, nil
}

// SavePlantFamily saves a new plant family.
func (s *Storage) SavePlantFamily(ctx context.Context, p *models.PlantFamily, recommendationId int64) (int64, error) {
	op := "storage.postgres.SavePlantFamily"
	_, err := s.db.Prepare("INSERT INTO plant_families (name, description, recommendation_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	// check if recommendation exists in recommendation table.
	row := s.db.QueryRowContext(ctx, "SELECT id FROM recommendations WHERE id = $1", recommendationId)

	err = row.Scan(&recommendationId)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := s.db.ExecContext(ctx, p.Name, p.Description, recommendationId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrPlantFamilyExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	// Return the ID of the inserted row.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// PlantFamily returns a plant family by ID.
func (s *Storage) PlantFamily(ctx context.Context, id int64) (models.PlantFamily, error) {
	op := "storage.postgres.PlantFamily"
	row := s.db.QueryRowContext(ctx, "SELECT id, name, description, recommendation_id FROM plant_families WHERE id = $1", id)
	var p models.PlantFamily
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.RecommendationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, storage.ErrPlantFamilyNotFound
		}
		return p, fmt.Errorf("%s: %v", op, err)
	}
	return p, nil
}

// PlantFamilyByName PlantFamily returns a plant family by name.
func (s *Storage) PlantFamilyByName(ctx context.Context, name string) (models.PlantFamily, error) {
	op := "storage.postgres.PlantFamilyByName"
	row := s.db.QueryRowContext(ctx, "SELECT id, name, description, recommendation_id FROM plant_families WHERE name = $1", name)
	var p models.PlantFamily
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.RecommendationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, storage.ErrPlantFamilyNotFound
		}
		return p, fmt.Errorf("%s: %v", op, err)
	}
	return p, nil
}
