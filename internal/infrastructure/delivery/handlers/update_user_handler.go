package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
)

func UpdateUserHandler(updateUserUseCase usecase.UpdateUserUseCase) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var user entities.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userEmail, ok := ctx.Get("email")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can not get user id"})
		}
		user, err := updateUserUseCase.Exec(userEmail.(string), user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"user": user})
	}
}
