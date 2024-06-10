package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"fmt"
	"time"
)

type ModuleDataPostgres struct {
	DB *sql.DB
}

func NewModuleDataPostgres(db *sql.DB) *ModuleDataPostgres {
	return &ModuleDataPostgres{DB: db}
}

func (s *ModuleDataPostgres) GetModuleData(moduleId int64, startDate, endDate time.Time) ([]models.ModuleData, error) {
	op := "storage.ModuleData.GetModuleData"
	var moduleData models.ModuleData

	rows, err := s.DB.Query("SELECT module_id, temperature, humidity_in, humidity_out, illuminance, measurement_time FROM moduleData WHERE module_id = $1 AND date >= $2 AND date <= $3", moduleId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	var moduleDataList []models.ModuleData
	for rows.Next() {
		err := rows.Scan(&moduleData.ModuleId, &moduleData.Temperature, &moduleData.HumidityIn, &moduleData.HumidityOut, &moduleData.Illuminance, &moduleData.MeasurementTime)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		moduleDataList = append(moduleDataList, moduleData)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return moduleDataList, nil

}

func (s *ModuleDataPostgres) SaveModuleData(moduleData models.ModuleData) error {
	op := "storage.ModuleData.SaveModuleData"
	_, err := s.DB.Exec("INSERT INTO moduleData (module_id, temperature, humidity_in, humidity_out, illuminance, measurement_time) VALUES ($1, $2, $3, $4, $5, $6)", moduleData.ModuleId, moduleData.Temperature, moduleData.HumidityIn, moduleData.HumidityOut, moduleData.Illuminance, moduleData.MeasurementTime)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
