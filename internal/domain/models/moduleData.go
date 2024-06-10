package models

import "time"

type ModuleData struct {
	ModuleId        int       `json:"module_id"`
	HumidityIn      float64   `json:"humidity_in"`
	HumidityOut     float64   `json:"humidity_out"`
	Temperature     float64   `json:"temperature"`
	Illuminance     float64   `json:"illuminance"`
	MeasurementTime time.Time `json:"measurement_time"`
}
