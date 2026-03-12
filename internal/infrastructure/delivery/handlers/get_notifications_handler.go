package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func GetNotificationsByUserIDHandler(GetNotificationByUserUseCase usecase.GetNotificationsByUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
			return
		}
		notifications, err := GetNotificationByUserUseCase.Exec(userID.(uint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]dtos.NotificationResponseDTO, 0)
		for _, notification := range notifications {
			response = append(response, dtos.NotificationToDto(notification))
		}

		ctx.JSON(http.StatusOK, gin.H{"notifications": response})
	}
}
