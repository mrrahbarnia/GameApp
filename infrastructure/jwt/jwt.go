package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO - Read this constant from .env later
const (
	secret = "test_secret"
)

type JWT struct {
	secret []byte
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func New() *JWT {
	return &JWT{
		secret: []byte(secret),
	}
}

func (j *JWT) GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) ValidateToken(tokenString string) (uint, error) {
	return 0, nil
}
