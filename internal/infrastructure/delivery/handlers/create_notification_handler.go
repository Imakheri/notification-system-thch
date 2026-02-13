package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
)

func CreateNotificationHandler(createNotificationUseCase usecase.CreateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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
		notification, err := createNotificationUseCase.Exec(userID.(uint), userEmail.(string), notification)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"notification": notification})
	}
}
