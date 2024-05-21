package jwt

import (
	"EMP_Back/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	salt := "EMP_Back_Token"
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["user_email"] = user.Email
	claims["user_role"] = user.IsAdmin
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
