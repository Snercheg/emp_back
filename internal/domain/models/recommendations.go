package models

type Recommendation struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	TemperatureMin  float64 `json:"temperature_min"`
	TemperatureMax  float64 `json:"temperature_max"`
	HumidityInMin   float64 `json:"humidity_in_min"`
	HumidityInMax   float64 `json:"humidity_in_max"`
	HumidityOutMin  float64 `json:"humidity_out_min"`
	HumidityOutMax  float64 `json:"humidity_out_max"`
	IlluminationMin float64 `json:"illuminance_min"`
	IlluminationMax float64 `json:"illuminance_max"`
	Description     string  `json:"description"`
}
