package models

import "errors"

type Module struct {
	ID            int          `json:"id"`
	Name          string       `json:"name" binding:"required"`
	PlantFamilyID int          `json:"plant_family_id"`
	PlantFamily   *PlantFamily `json:"plant_family"`
	Status        string       `json:"status"`
}

type UpdateModuleInput struct {
	Name          string `json:"name"`
	PlantFamilyID int    `json:"plant_family_id"`
}

func (i UpdateModuleInput) Validate() error {
	if i.Name == "" && i.PlantFamilyID == 0 {
		return errors.New("update structure has no values")
	}
	return nil
}
