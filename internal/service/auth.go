package service

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/jwt"
	"EMP_Back/internal/storage"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	repo     storage.Repository
	log      *slog.Logger
	tokenTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewAuthService(
	repo storage.Repository,
	log *slog.Logger,
) *AuthService {
	return &AuthService{
		log:  log,
		repo: repo}
}

func (s *AuthService) SaveUser(user models.User) (int64, error) {
	const op = "app.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
	)

	log.Info("Registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.PassHash), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	user.PassHash = string(passHash)
	log.Info("User password hash generated")

	return s.repo.SaveUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	const op = "app.GenerateToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Generating token")
	user, err := s.repo.GetUser(email)
	if err != nil {
		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		s.log.Warn("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	token, err := jwt.NewToken(user, s.tokenTTL)
	if err != nil {
		s.log.Error("failed to create token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
