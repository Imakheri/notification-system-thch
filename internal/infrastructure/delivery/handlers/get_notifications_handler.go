package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Security BearerAuth
// @Summary List a specific user's notifications
// @Description Gets all notifications linked to the authenticated user's ID
// @Tags notifications
// @Produce json
// @Success 200  {object}  []dtos.NotificationResponseDTO "List of notifications"
// @Failure 400  {object}  dtos.ErrorResponseDTO "can not get user id"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /notifications [get]
func GetNotificationsByUserIDHandler(GetNotificationByUserUseCase usecase.GetNotificationsByUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user id"})
			return
		}
		notifications, err := GetNotificationByUserUseCase.Exec(userID.(uint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		response := make([]dtos.NotificationResponseDTO, 0)
		for _, notification := range notifications {
			response = append(response, dtos.NotificationToDto(notification))
		}

		ctx.JSON(http.StatusOK, gin.H{"notifications": response})
	}
}
