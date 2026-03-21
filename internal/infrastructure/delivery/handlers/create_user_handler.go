package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Summary Create a new user
// @Description Register a new user in the system by providing their basic information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.CreateUserDTO true "New user data"
// @Success 201  {object}  dtos.ResponseUserDTO "user created successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "invalid user properties"
// @Failure 404  {object}  dtos.ErrorResponseDTO "record not found"
// @Failure 500  {object}  dtos.ErrorResponseDTO "user already exists"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /users [post]
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
