package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/infrastructure/service"
)

func AuthorizeJWT(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := ctx.GetHeader("Authorization")

		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "must enter a valid authorization header"})
		}

		tokenString := authHeader[len(BEARER_SCHEMA)+1:]

		token, err := service.NewJWTService(cfg).ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("id", uint(claims["id"].(float64)))
			ctx.Set("email", claims["email"].(string))
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
