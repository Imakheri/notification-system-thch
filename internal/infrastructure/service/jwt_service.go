package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/joho/godotenv"
)

type JWTService struct {
	jwToken string
}

func NewJWTService() gateway.JwTokenService {
	return &JWTService{}
}

func (s JWTService) GenerateToken(email string, id uint) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": email,
		"id":    id,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 60).Unix(),
	})

	token, err := claims.SignedString([]byte(getSecretKey("SECRET_JWT")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecretKey("SECRET_JWT")), nil
	})
}

func getSecretKey(variableName string) string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	secret := os.Getenv(variableName)
	return secret
}
