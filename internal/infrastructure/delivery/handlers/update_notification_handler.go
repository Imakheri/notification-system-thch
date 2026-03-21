package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Security BearerAuth
// @Summary Update a notification
// @Description Modify the content of a specific notification by verifying that it belongs to the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Param notification body dtos.UpdateNotificationDTO true "Notification information for update"
// @Success 200  {object}  dtos.ResponseUserDTO "notification updated created successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "must enter a valid id"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "recipient does not exist"
// @Failure 500  {object}  dtos.ErrorResponseDTO "invalid recipient"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /notification/{id} [put]
func UpdateNotificationHandler(updateNotificationUseCase usecase.UpdateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.UpdateNotificationDTO
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		idParam := ctx.Param("id")
		if len(idParam) <= 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: "must enter a valid id"})
			return
		}
		notificationID, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		userID, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user id"})
			return
		}

		notification, err := dtos.NotificationUpdateToEntity(input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		updatedNotification, err := updateNotificationUseCase.Exec(userID.(uint), notificationID, notification)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"notification": dtos.NotificationToDto(updatedNotification)})
	}
}
