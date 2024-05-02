package models

type PlantFamily struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	RecommendationId int             `json:"recommendation_id"`
	Recommendation   *Recommendation `json:"recommendation"`
}
