package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func UpdateNotificationHandler(updateNotificationUseCase usecase.UpdateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.UpdateNotificationDTO
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idParam := ctx.Param("id")
		if len(idParam) <= 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "must enter a valid id"})
		}
		notificationID, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		userID, ok := ctx.Get("id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
		}

		notification, err := dtos.NotificationUpdateToEntity(input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		updatedNotification, err := updateNotificationUseCase.Exec(userID.(uint), notificationID, notification)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"notification": dtos.NotificationToDto(updatedNotification)})
	}
}
