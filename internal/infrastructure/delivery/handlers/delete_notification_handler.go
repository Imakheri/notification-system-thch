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

		id, err := strconv.Atoi(idParam)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		notificationID, err := deleteNotificationUseCase.Exec(id)
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
