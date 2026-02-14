package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ValidateAPIKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
		secret := os.Getenv("API_KEY")
		apiKey := ctx.GetHeader("api_key")

		if len(apiKey) == 0 || apiKey != secret {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "must enter a valid API Key"})
		}
	}
}
