package models

type Setting struct {
	ModuleId        int     `json:"module_id"`
	Name            string  `json:"name"`
	TemperatureMin  float64 `json:"temperature_min"`
	TemperatureMax  float64 `json:"temperature_max"`
	HumidityInMin   float64 `json:"humidity_in_min"`
	HumidityInMax   float64 `json:"humidity_in_max"`
	HumidityOutMin  float64 `json:"humidity_out_min"`
	HumidityOutMax  float64 `json:"humidity_out_max"`
	IlluminationMin float64 `json:"illumination_min"`
	IlluminationMax float64 `json:"illumination_max"`
}
