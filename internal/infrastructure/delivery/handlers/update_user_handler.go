package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func UpdateUserHandler(updateUserUseCase usecase.UpdateUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.UpdateUserDTO
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userEmail, ok := ctx.Get("email")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can not get user email"})
		}
		updateUserEntity, err := dtos.UserUpdateToEntity(input)
		updatedUser, err := updateUserUseCase.Exec(userEmail.(string), updateUserEntity)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"user": dtos.UserResponseToDTO(updatedUser)})
	}
}
