package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func GetUserHandler(useCase usecase.GetUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.GetUserDTO
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userEntity := dtos.GetUserToEntity(input)
		user, token, err := useCase.Exec(userEntity)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user":  dtos.UserResponseToDTO(user),
			"token": token,
		})
	}
}
