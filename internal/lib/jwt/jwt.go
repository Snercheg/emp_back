package jwt

import (
	"EMP_Back/internal/domain/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	signingKey = "EMP_Back_Token"
	salt       = "EMP_Back_Token"
)

type tokenClaims struct {
	jwt.Claims
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  bool   `json:"user_role"`
}

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &tokenClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		UserRole:  user.IsAdmin,
	}
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(accessToken string) (int64, error) {
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	if token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("token is invalid")
}
