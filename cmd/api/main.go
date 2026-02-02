package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers"
	"github.com/imakheri/notifications-thch/internal/infrastructure/middleware"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
)

var prefix = "/api/v1/"

func main() {
	router := gin.Default()
	db := repository.Database()
	repository.Migration()
	userReposiroty := repository.NewUserRepository(db)
	createUseCase := usecase.NewCreateUser(userReposiroty)

	router.POST(prefix+"/users", handlers.CreateUserHandler(createUseCase))
	getUserUseCase := usecase.NewGetUser(userReposiroty)
	router.POST(prefix+"/user", handlers.GetUserHandler(getUserUseCase))

	//CREATE
	notificationRepo := repository.NewNotificationRepository(db)
	createNotificationUseCase := usecase.NewCreateNotificationUseCase(notificationRepo)
	router.POST(prefix+"/notification", middleware.AuthorizeJWT(), handlers.CreateNotificationHandler(createNotificationUseCase))

	getNotificationsByUserUseCase := usecase.NewGetNotificationsByUserUseCase(notificationRepo)
	//READ
	router.GET(prefix+"/notifications", middleware.AuthorizeJWT(), handlers.GetNotificationsByUserIDHandler(getNotificationsByUserUseCase))
	//UPDATE
	updateNotificationUseCase := usecase.NewUpdateNotificationUseCase(notificationRepo)
	router.PUT(prefix+"/notification/:id", middleware.AuthorizeJWT(), handlers.UpdateNotificationHandler(updateNotificationUseCase))
	//DELETE
	deleteNotificationUsecase := usecase.NewDeleteNotification(notificationRepo)
	router.DELETE(prefix+"/notification/:id", middleware.AuthorizeJWT(), handlers.DeleteNotificationHandler(deleteNotificationUsecase))

	router.Run()
}
