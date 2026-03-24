package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Security BearerAuth
// @Summary Create a new notification
// @Description Create a notification linked to the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body dtos.CreateNotificationDTO true "new notification data"
// @Success 200  {object}  dtos.ResponseUserDTO "notification created successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "invalid notification properties"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "can not get user id"
// @Failure 500  {object}  dtos.ErrorResponseDTO "can not get user email"
// @Failure 500  {object}  dtos.ErrorResponseDTO "recipient does not exist"
// @Failure 500  {object}  dtos.ErrorResponseDTO "invalid recipient"
// @Failure 500  {object}  dtos.ErrorResponseDTO "cannot get sender information"
// @Failure 500  {object}  dtos.ErrorResponseDTO "an error occurred while trying to send"
// @Failure 500  {object}  dtos.ErrorResponseDTO "failed to encrypt password"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /notification [post]
func CreateNotificationHandler(createNotificationUseCase usecase.CreateNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var notification dtos.CreateNotificationDTO
		if err := ctx.ShouldBind(&notification); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		notificationEntity, err := dtos.NotificationToEntity(notification)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		userID, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user id"})
			return
		}
		userEmail, ok := ctx.Get("email")
		if !ok || userEmail == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user email"})
			return
		}

		newNotification, err := createNotificationUseCase.Exec(userID.(uint), userEmail.(string), notificationEntity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"notification": dtos.NotificationToDto(newNotification)})
	}
}
