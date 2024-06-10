package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type SettingsPostgres struct {
	DB *sql.DB
}

func NewSettingsPostgres(db *sql.DB) *SettingsPostgres {
	return &SettingsPostgres{DB: db}
}

func (s *SettingsPostgres) SaveSetting(moduleId int64, settings models.Setting) (int64, error) {
	op := "storage.Settings.SaveSetting"
	result, err := s.DB.Exec("INSERT INTO settings (module_id, name, temperature_min, temperature_max, humidity_in_min, humidity_in_max, humidity_out_min, humidity_out_max, illumination_min, illumination_max) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", moduleId, settings.Name, settings.TemperatureMin, settings.TemperatureMax, settings.HumidityInMin, settings.HumidityInMax, settings.HumidityOutMin, settings.HumidityOutMax, settings.IlluminationMin, settings.IlluminationMax)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, ErrSettingExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

func (s *SettingsPostgres) GetSetting(moduleId int64) (*models.Setting, error) {
	op := "storage.Settings.GetSetting"
	var setting models.Setting
	err := s.DB.QueryRow("SELECT * FROM settings WHERE module_id = $1", moduleId).Scan(&setting.ModuleId, &setting.ModuleId, &setting.Name, &setting.TemperatureMin, &setting.TemperatureMax, &setting.HumidityInMin, &setting.HumidityInMax, &setting.HumidityOutMin, &setting.HumidityOutMax, &setting.IlluminationMin, &setting.IlluminationMax)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSettingNotFound
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return &setting, nil
}

func (s *SettingsPostgres) UpdateSetting(moduleId int64, setting models.Setting) error {
	op := "storage.Settings.UpdateSetting"
	_, err := s.DB.Exec("UPDATE settings SET name = $1, temperature_min = $2, temperature_max = $3, humidity_in_min = $4, humidity_in_max = $5, humidity_out_min = $6, humidity_out_max = $7, illumination_min = $8, illumination_max = $9 WHERE module_id = $1", moduleId, setting.Name, setting.TemperatureMin, setting.TemperatureMax, setting.HumidityInMin, setting.HumidityInMax, setting.HumidityOutMin, setting.HumidityOutMax, setting.IlluminationMin, setting.IlluminationMax, moduleId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23504" { // not-null violation
			return fmt.Errorf("%s: %v", op, err, ErrSettingCannotBeChanged)
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
