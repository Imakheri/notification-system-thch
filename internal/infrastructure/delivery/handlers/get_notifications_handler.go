package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
)

func GetNotificationsByUserIDHandler(GetNotificationByUserUseCase usecase.GetNotificationsByUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
		}
		notifications, err := GetNotificationByUserUseCase.Exec(userID.(uint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"notifications": notifications})
	}
}
