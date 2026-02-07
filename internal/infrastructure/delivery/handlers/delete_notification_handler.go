package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
)

func DeleteNotificationHandler(deleteNotificationUseCase usecase.DeleteNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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

		notificationID, err = deleteNotificationUseCase.Exec(userID.(uint), notificationID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification deleted successfully",
			"id":      notificationID,
		})
	}
}
