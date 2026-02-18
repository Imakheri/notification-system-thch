package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
)

type APIKey struct {
	config config.Config
}

func NewAPIKey(cfg *config.Config) *APIKey {
	return &APIKey{
		config: *cfg,
	}
}

func (a *APIKey) ValidateAPIKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("api_key")
		if len(apiKey) <= 0 || apiKey != a.config.ApiKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "must enter a valid API Key"})
		}
		ctx.Next()
	}
}
