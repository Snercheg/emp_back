package models

type setting struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	TemperatureMin  float64 `json:"temperature_min"`
	TemperatureMax  float64 `json:"temperature_max"`
	HumidityMin     float64 `json:"humidity_min"`
	HumidityMax     float64 `json:"humidity_max"`
	IlluminationMin float64 `json:"illumination_min"`
	IlluminationMax float64 `json:"illumination_max"`
}
