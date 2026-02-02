package gateway

import "github.com/golang-jwt/jwt/v5"

type JwTokenService interface {
	GenerateToken(email string, id uint) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}
