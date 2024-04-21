package models

import "time"

type Data struct {
	ModuleId        int           `json:"module_id"`
	Humidity        float64       `json:"humidity"`
	Temperature     float64       `json:"temperature"`
	Illuminance     float64       `json:"illuminance"`
	MeasurementTime time.Duration `json:"measurement_time"`
}
