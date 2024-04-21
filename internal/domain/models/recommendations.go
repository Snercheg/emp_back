package models

type Recommendation struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	TemperatureMin float64 `json:"temperature_min"`
	TemperatureMax float64 `json:"temperature_max"`
	HumidityMin    float64 `json:"humidity_min"`
	HumidityMax    float64 `json:"humidity_max"`
	IlluminanceMin float64 `json:"illuminance_min"`
	IlluminanceMax float64 `json:"illuminance_max"`
	Description    string  `json:"description"`
	//CreatedAt      time.Duration  `json:"created_at"`
}
