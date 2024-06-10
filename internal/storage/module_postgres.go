package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strconv"
	"strings"
)

type ModulePostgres struct {
	DB *sql.DB
}

func NewModulePostgres(db *sql.DB) *ModulePostgres {
	return &ModulePostgres{DB: db}
}

// SaveModule saves a module in the database
func (m *ModulePostgres) SaveModule(userId int64, module models.Module) (int64, error) {
	op := "storage.repository.SaveUser"

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	row := tx.QueryRow("SELECT * FROM recommendations WHERE id = (SELECT recommendation_id FROM plantFamily WHERE id = $1)", module.PlantFamilyID)
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
	err = tx.QueryRow("INSERT INTO modules (user_id, plant_family_id, name) VALUES ($1, $2, $3) RETURNING id", userId, module.PlantFamilyID, module.Name).Scan(&id)
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
	var modules []models.Module
	for rows.Next() {
		var moduleID int64
		err = rows.Scan(&moduleID)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		// Get module details from modules table using module_id
		module, err := m.GetModule(userID, moduleID)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		modules = append(modules, *module)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (m *ModulePostgres) GetModule(userId, moduleId int64) (*models.Module, error) {
	op := "storage.repository.GetModule"
	var module models.Module
	err := m.DB.QueryRow("SELECT m.id, m.name FROM modules m JOIN user_modules um ON m.id = um.module_id WHERE m.id = $1 AND um.user_id = $2", moduleId, userId).Scan(&module.ID, &module.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModuleNotFound
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return &module, nil
}

func (m *ModulePostgres) DeleteModule(userId, moduleId int64) error {
	op := "storage.repository.DeleteModule"
	_, err := m.DB.Exec("DELETE FROM modules m USING user_modules um WHERE m.id = $1 AND um.module_id = m.id AND um.user_id = $2", moduleId, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key
			return fmt.Errorf("%s: %v", op, err, ErrModuleNotFound)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (m *ModulePostgres) UpdateModule(userId, moduleId int64, module models.UpdateModuleInput) error {
	op := "storage.repository.UpdateModuleInput"
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if module.Name != "" {
		setValues = append(setValues, "name = $"+strconv.Itoa(argId))
		argId++
		args = append(args, module.Name)
	}
	if module.PlantFamilyID != 0 {
		setValues = append(setValues, "plant_family_id = $"+strconv.Itoa(argId))
		argId++
		args = append(args, module.PlantFamilyID)
	}
	if len(setValues) > 0 {
		_, err := m.DB.Exec("UPDATE modules m SET "+strings.Join(setValues, ", ")+" FROM user_modules um WHERE m.id = $1 AND um.module_id = m.id AND um.user_id = $2", moduleId, userId, module.Name, module.PlantFamilyID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
				return fmt.Errorf("%s: %v", op, ErrModuleCannotBeChanged)
			}
			return fmt.Errorf("%s: %v", op, err)
		}
	}
	return nil
}
