package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
)

// @Security ApiKeyAuth
// @Summary Loggin / Get user data
// @Description Authenticates an user and return thier data along a JWToken
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.GetUserDTO true "Access credentials"
// @Success 201  {object}  dtos.LoginResponseDTO "user logged in successfully"
// @Failure 400  {object}  dtos.ErrorResponseDTO "invalid user credentials"
// @Failure 401  {object}  dtos.ErrorResponseDTO "the e-mail address or password is incorrect"
// @Failure 500  {object}  dtos.ErrorResponseDTO "could not generate JWT"
// @Failure 500  {object}  dtos.ErrorResponseDTO "internal server error"
// @Router /user [post]
func GetUserHandler(getUserUseCase usecase.GetUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var input dtos.GetUserDTO
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}
		userEntity := dtos.GetUserToEntity(input)
		user, token, err := getUserUseCase.Exec(userEntity)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponseDTO{Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, dtos.LoginResponseDTO{
			User:  dtos.UserResponseToDTO(user),
			Token: token,
		})
	}
}
