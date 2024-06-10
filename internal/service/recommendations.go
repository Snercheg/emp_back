package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"log/slog"
)

type RecommendationsService struct {
	repo storage.Recommendations
	log  *slog.Logger
}

func NewRecommendationsService(repo storage.Recommendations, log *slog.Logger) *RecommendationsService {
	return &RecommendationsService{repo: repo, log: log}
}

func (s *RecommendationsService) GetRecommendation(id int64) (models.Recommendation, error) {
	return s.repo.GetRecommendation(id)
}
func (s *RecommendationsService) GetRecommendations() ([]models.Recommendation, error) {
	return s.repo.GetRecommendations()
}

func (s *RecommendationsService) SaveRecommendation(recommendation models.Recommendation) (int64, error) {
	return s.repo.SaveRecommendation(recommendation)
}

func (s *RecommendationsService) UpdateRecommendation(id int64, recommendation models.Recommendation) error {
	return s.repo.UpdateRecommendation(id, recommendation)
}

func (s *RecommendationsService) DeleteRecommendation(id int64) error {
	return s.repo.DeleteRecommendation(id)
}
