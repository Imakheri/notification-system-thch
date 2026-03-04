package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func CreateNotificationHandler(createNotificationUseCase usecase.CreateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var notification dtos.CreateNotificationDTO
		if err := ctx.ShouldBind(&notification); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		notificationEntity, err := dtos.NotificationToEntity(notification)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		userID, ok := ctx.Get("id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
		}
		userEmail, ok := ctx.Get("email")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can not get user email"})
		}

		newNotification, err := createNotificationUseCase.Exec(userID.(uint), userEmail.(string), notificationEntity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"notification": dtos.NotificationToDto(newNotification)})
	}
}
