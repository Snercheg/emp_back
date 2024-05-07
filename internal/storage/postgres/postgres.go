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

	stmt, err := s.db.Prepare("INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.ExecContext(ctx, username, email, passHash)
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
	stmt, err := s.db.Prepare("INSERT INTO recommendations (title, temperature_min, temperature_max, humidity_min, humidity_max, illumination_min, illumination_max, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	res, err := stmt.ExecContext(ctx, r.Title, r.TemperatureMin, r.TemperatureMax, r.HumidityMin, r.HumidityMax, r.IlluminationMin, r.IlluminationMax, r.Description)
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

// UpdateRecommendation updates a recommendation.
func (s *Storage) UpdateRecommendation(ctx context.Context, r *models.Recommendation) error {
	op := "storage.postgres.UpdateRecommendation"
	stmt, err := s.db.Prepare("UPDATE recommendations SET title = $1, temperature_min = $2, temperature_max = $3, humidity_min = $4, humidity_max = $5, illumination_min = $6, illumination_max = $7, description = $8 WHERE id = $9")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.ExecContext(ctx, r.Title, r.TemperatureMin, r.TemperatureMax, r.HumidityMin, r.HumidityMax, r.IlluminationMin, r.IlluminationMax, r.Description, r.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, storage.ErrRecommendationCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

// SavePlantFamily saves a new plant family.
func (s *Storage) SavePlantFamily(ctx context.Context, p *models.PlantFamily, recommendationId int64) (int64, error) {
	op := "storage.postgres.SavePlantFamily"

	// check if recommendation exists in recommendation table.
	row := s.db.QueryRowContext(ctx, "SELECT id FROM recommendations WHERE id = $1", recommendationId)

	err := row.Scan(&recommendationId)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	stmt, err := s.db.Prepare("INSERT INTO plant_families (name, description, recommendation_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.ExecContext(ctx, p.Name, p.Description, recommendationId)
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

// UpdatePlantFamily updates a plant family.
func (s *Storage) UpdatePlantFamily(ctx context.Context, p *models.PlantFamily) error {
	op := "storage.postgres.UpdatePlantFamily"
	stmt, err := s.db.Prepare("UPDATE plant_families SET name = $1, description = $2, recommendation_id = $3 WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.ExecContext(ctx, p.Name, p.Description, p.RecommendationId, p.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, storage.ErrPlantFamilyCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
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

// SaveModule saves a new module.
func (s *Storage) SaveModule(ctx context.Context, m *models.Module) (int64, error) {
	op := "storage.postgres.SaveModule"
	// check if setting exists in setting table.
	row := s.db.QueryRowContext(ctx, "SELECT id FROM settings WHERE id = $1", m.SettingID)
	var settingId int64
	err := row.Scan(&settingId)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	// check if plant family exists in plant_families table.
	row = s.db.QueryRowContext(ctx, "SELECT id FROM plant_families WHERE id = $1", m.PlantFamilyID)
	var plantFamilyId int64
	err = row.Scan(&plantFamilyId)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	stmt, err := s.db.Prepare("INSERT INTO modules (name, setting_id, plant_family_id ) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	res, err := stmt.ExecContext(ctx, m.Name, settingId, plantFamilyId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrModuleExist)
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

// Module returns a module by ID.
func (s *Storage) Module(ctx context.Context, id int64) (models.Module, error) {
	op := "storage.postgres.Module"
	row := s.db.QueryRowContext(ctx, "SELECT id, name, setting_id, plant_family_id FROM modules WHERE id = $1", id)
	var m models.Module
	err := row.Scan(&m.ID, &m.Name, &m.SettingID, &m.PlantFamilyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m, storage.ErrModuleNotFound
		}
		return m, fmt.Errorf("%s: %v", op, err)
	}
	return m, nil
}

// SaveUserModule saves a new user module.
func (s *Storage) SaveUserModule(ctx context.Context, m *models.UserModule) (int64, error) {
	op := "storage.postgres.SaveUserModule"
	stmt, err := s.db.Prepare("INSERT INTO user_modules (user_id, module_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.ExecContext(ctx, m.UserID, m.ModuleID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrUserModuleExist)
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

// UserModule returns a user module by ID.
func (s *Storage) UserModule(ctx context.Context, id int64) (models.UserModule, error) {
	op := "storage.postgres.UserModule"
	row := s.db.QueryRowContext(ctx, "SELECT id, user_id, module_id FROM user_modules WHERE id = $1", id)
	var m models.UserModule
	err := row.Scan(&m.ID, &m.UserID, &m.ModuleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m, storage.ErrUserModuleNotFound
		}
		return m, fmt.Errorf("%s: %v", op, err)
	}
	return m, nil
}

// DeleteUserModule deletes a user module by ID.
func (s *Storage) DeleteUserModule(ctx context.Context, id int64) error {
	op := "storage.postgres.DeleteUserModule"
	stmt, err := s.db.Prepare("DELETE FROM user_modules WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key
			return fmt.Errorf("%s: %v", op, err, storage.ErrUserModuleNotFound)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

// SaveSetting saves a new setting.
func (s *Storage) SaveSetting(ctx context.Context, set *models.Setting) (int64, error) {
	op := "storage.postgres.SaveSetting"
	stmt, err := s.db.Prepare("INSERT INTO settings (name, temperature_min, temperature_max, humidity_min, humidity_max, illumination_min, illumination_max) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	res, err := stmt.ExecContext(ctx, set.Name, set.TemperatureMin, set.TemperatureMax, set.HumidityMin, set.HumidityMax, set.IlluminationMin, set.IlluminationMax)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrSettingExist)
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

// Setting returns a setting by ID.
func (s *Storage) Setting(ctx context.Context, id int64) (models.Setting, error) {
	op := "storage.postgres.Setting"
	row := s.db.QueryRowContext(ctx, "SELECT id, name, temperature_min, temperature_max, humidity_min, humidity_max, illumination_min, illumination_max FROM settings WHERE id = $1", id)
	var set models.Setting
	err := row.Scan(&set.ID, &set.Name, &set.TemperatureMin, &set.TemperatureMax, &set.HumidityMin, &set.HumidityMax, &set.IlluminationMin, &set.IlluminationMax)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return set, storage.ErrSettingNotFound
		}
		return set, fmt.Errorf("%set: %v", op, err)
	}
	return set, nil
}

// UpdateSetting updates a setting.
func (s *Storage) UpdateSetting(ctx context.Context, set *models.Setting) error {
	op := "storage.postgres.UpdateSetting"
	stmt, err := s.db.Prepare("UPDATE settings SET name = $1, temperature_min = $2, temperature_max = $3, humidity_min = $4, humidity_max = $5, illumination_min = $6, illumination_max = $7 WHERE id = $8")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	_, err = stmt.ExecContext(ctx, set.Name, set.TemperatureMin, set.TemperatureMax, set.HumidityMin, set.HumidityMax, set.IlluminationMin, set.IlluminationMax, set.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not null
			return fmt.Errorf("%s: %v", op, storage.ErrSettingCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
