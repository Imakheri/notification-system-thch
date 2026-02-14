package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type JWTService struct {
	jwToken string
	config  *config.Config
}

func NewJWTService(config *config.Config) gateway.JwTokenService {
	return &JWTService{
		config: config,
	}
}

func (s JWTService) GenerateToken(email string, id uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": email,
		"id":    id,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 60).Unix(),
	})

	token, err := claims.SignedString([]byte(s.config.SecretJWT))
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
		return []byte(s.config.SecretJWT), nil
	})
}
