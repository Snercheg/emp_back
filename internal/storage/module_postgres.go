package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type ModulePostgres struct {
	DB *sql.DB
}

func NewModulePostgres(db *sql.DB) *ModulePostgres {
	return &ModulePostgres{DB: db}
}

// SaveModule saves a module in the database
func (m *ModulePostgres) SaveModule(userId, plantFamilyId int64, module *models.Module) (int64, error) {
	op := "storage.repository.SaveUser"

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	row := tx.QueryRow("SELECT * FROM recommendations WHERE id = (SELECT recommendation_id FROM plantFamily WHERE id = $1)", plantFamilyId)
	var recommendation models.Recommendation
	err = row.Scan(&recommendation.TemperatureMin, &recommendation.TemperatureMax, &recommendation.HumidityInMin, &recommendation.HumidityInMax, &recommendation.HumidityOutMin, &recommendation.HumidityOutMax, &recommendation.IlluminationMin, &recommendation.IlluminationMax)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecommendationNotFound
		}
		err := tx.Rollback()
		if err != nil {
			return 0, fmt.Errorf("%s: %v", op, err)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	var id int64
	err = tx.QueryRow("INSERT INTO modules (user_id, plant_family_id, name) VALUES ($1, $2, $3) RETURNING id", userId, plantFamilyId, module.Name).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, ErrModuleExist)
		}
		err := tx.Rollback()
		if err != nil {
			return 0, fmt.Errorf("%s: %v", op, err)
		}
		return 0, fmt.Errorf("%s: %v", op, err)

	}
	err = tx.QueryRow("INSERT INTO settings (module_id, temperature_min, temperature_max, humidity_in_min, humidity_in_max, humidity_out_min, humidity_out_max, illumination_min, illumination_max) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", id, recommendation.TemperatureMin, recommendation.TemperatureMax, recommendation.HumidityInMin, recommendation.HumidityInMax, recommendation.HumidityOutMin, recommendation.HumidityOutMax, recommendation.IlluminationMin, recommendation.IlluminationMax).Scan()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, fmt.Errorf("%s: %v", op, err)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	err = tx.QueryRow("INSERT INTO user_modules (user_id, module_id) VALUES ($1, $2)", userId, id).Scan()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, fmt.Errorf("%s: %v", op, err)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// GetModules returns all modules of a plant family
func (m *ModulePostgres) GetModules(userID int64) ([]models.Module, error) {
	op := "storage.repository.GetModules"

	// Select module_id from user_modules table
	rows, err := m.DB.Query("SELECT module_id FROM user_modules WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var moduleID int64
		err = rows.Scan(&moduleID)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		// Get module details from modules table using module_id
		module, err := m.GetModule(moduleID)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		modules = append(modules, *module)
	}

	return modules, nil
}

func (m *ModulePostgres) GetModule(id int64) (*models.Module, error) {
	op := "storage.repository.GetModule"
	var module models.Module
	err := m.DB.QueryRow("SELECT id, name FROM modules WHERE id = $1", id).Scan(&module.ID, &module.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModuleNotFound
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return &module, nil
}

func (m *ModulePostgres) DeleteModule(id int64) error {
	op := "storage.repository.DeleteModule"
	_, err := m.DB.Exec("DELETE FROM modules WHERE id = $1", id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key
			return fmt.Errorf("%s: %v", op, err, ErrModuleNotFound)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (m *ModulePostgres) UpdateModule(id int64, module *models.Module) error {
	op := "storage.repository.UpdateModule"
	_, err := m.DB.Exec("UPDATE modules SET name = $1 WHERE id = $2", module.Name, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, ErrModuleCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
