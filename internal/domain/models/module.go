package models

type Module struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	SettingID     int          `json:"setting_id"`
	Setting       *Setting     `json:"setting"`
	PlantFamilyID int          `json:"plant_family_id"`
	PlantFamily   *PlantFamily `json:"plant_family"`
	Status        string       `json:"status"`
}
