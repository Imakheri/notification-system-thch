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
// @Summary Delete a specific notification
// @Description Delete a specific system notification by verifying that it belongs to the user
// @Tags notifications
// @Produce json
// @Param id   path      int  true  "Notification ID"
// @Success 200  {object}  dtos.SuccessfulDeleteResponseDTO "Notification deleted successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "must enter a valid id"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "can not get user id"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /notification/{id} [delete]
func DeleteNotificationHandler(deleteNotificationUseCase usecase.DeleteNotificationUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")

		if len(idParam) <= 0 || idParam == "0" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: "must enter a valid id"})
			return
		}

		notificationIDParam, err := strconv.Atoi(idParam)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		userID, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user id"})
			return
		}

		notificationID, err := deleteNotificationUseCase.Exec(userID.(uint), uint(notificationIDParam))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, dtos.SuccessfulDeleteResponseDTO{
			ID:      notificationID,
			Message: "Notification deleted successfully",
		})
	}
}
