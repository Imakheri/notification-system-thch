package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers"
	"github.com/imakheri/notifications-thch/internal/infrastructure/middleware"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
)

var prefix = "/api/v1/"

func main() {
	router := gin.Default()
	cfg := config.Load()
	db := repository.NewDatabase(cfg)

	userRepository := repository.NewUserRepository(db)
	createUseCase := usecase.NewCreateUser(userRepository)

	router.POST(prefix+"/users", middleware.NewAPIKey(cfg).ValidateAPIKey(), handlers.CreateUserHandler(createUseCase))
	getUserUseCase := usecase.NewGetUser(userRepository, cfg)
	router.POST(prefix+"/user", middleware.NewAPIKey(cfg).ValidateAPIKey(), handlers.GetUserHandler(getUserUseCase))
	//UPDATE USER
	//router.PUT(prefix+"/user/:id", middleware.ValidateAPIKey(), middleware.AuthorizeJWT(), handlers.UpdateNotificationHandler(updateNotificationUseCase))

	//CREATE
	notificationRepo := repository.NewNotificationRepository(db)
	createNotificationUseCase := usecase.NewCreateNotificationUseCase(notificationRepo, userRepository)
	router.POST(prefix+"/notification", middleware.NewAPIKey(cfg).ValidateAPIKey(), middleware.NewAuthorizeJWT().AuthorizeJWT(), handlers.CreateNotificationHandler(createNotificationUseCase))

	getNotificationsByUserUseCase := usecase.NewGetNotificationsByUserUseCase(notificationRepo)
	//READ
	router.GET(prefix+"/notifications", middleware.NewAPIKey(cfg).ValidateAPIKey(), middleware.NewAuthorizeJWT().AuthorizeJWT(), handlers.GetNotificationsByUserIDHandler(getNotificationsByUserUseCase))
	//UPDATE
	updateNotificationUseCase := usecase.NewUpdateNotificationUseCase(notificationRepo, userRepository)
	router.PUT(prefix+"/notification/:id", middleware.NewAPIKey(cfg).ValidateAPIKey(), middleware.NewAuthorizeJWT().AuthorizeJWT(), handlers.UpdateNotificationHandler(updateNotificationUseCase))
	//DELETE
	deleteNotificationUsecase := usecase.NewDeleteNotification(notificationRepo)
	router.DELETE(prefix+"/notification/:id", middleware.NewAPIKey(cfg).ValidateAPIKey(), middleware.NewAuthorizeJWT().AuthorizeJWT(), handlers.DeleteNotificationHandler(deleteNotificationUsecase))

	router.Run()
}
