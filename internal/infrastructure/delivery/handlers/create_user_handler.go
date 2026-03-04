package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

func CreateUserHandler(createUserUseCase usecase.CreateUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.CreateUserDTO
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userEntity, err := dtos.CreateUserToEntity(input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := createUserUseCase.Exec(userEntity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"user": dtos.UserResponseToDTO(newUser)})
	}
}
