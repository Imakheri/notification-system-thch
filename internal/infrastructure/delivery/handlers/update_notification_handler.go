package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
)

func UpdateNotificationHandler(updateNotificationUseCase usecase.UpdateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")

		if len(idParam) <= 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "must enter a valid id"})
		}

		notificationID, err := strconv.Atoi(idParam)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var notification entities.Notification
		if err := ctx.ShouldBind(&notification); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, ok := ctx.Get("id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
		}
		userEmail, ok := ctx.Get("email")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can not get user email"})
		}
		notification, err = updateNotificationUseCase.Exec(userID.(uint), userEmail.(string), notificationID, notification)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"notification": notification})
	}
}
