package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Security BearerAuth
// @Summary Update user data
// @Description Allows the authenticated user to update their information (except email)
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.UpdateUserDTO true "User information for update"
// @Success 200  {object}  dtos.ResponseUserDTO "user updated successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "invalid user properties"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "can not get user email"
// @Failure 500  {object}  dtos.ErrorResponseDTO "failed to encrypt password"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /user [put]
func UpdateUserHandler(updateUserUseCase usecase.UpdateUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.UpdateUserDTO
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}
		userEmail, ok := ctx.Get("email")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: "can not get user email"})
			return
		}
		updateUserEntity, err := dtos.UserUpdateToEntity(input)
		updatedUser, err := updateUserUseCase.Exec(userEmail.(string), updateUserEntity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": dtos.UserResponseToDTO(updatedUser)})
	}
}
