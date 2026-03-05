package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/routes"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/service"
)

func main() {
	router := gin.Default()
	cfg := config.Load()
	db := repository.NewDatabase(cfg)

	userRepository := repository.NewUserRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	channelRepository := repository.NewChannelRepository(db)

	simulatedApiService := service.NewSimulatedApiService()

	createUserUseCase := usecase.NewCreateUser(userRepository)
	getUserUseCase := usecase.NewGetUser(userRepository, cfg)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepository)

	createNotificationUseCase := usecase.NewCreateNotificationUseCase(notificationRepository, userRepository, channelRepository, simulatedApiService)
	getNotificationsByUserUseCase := usecase.NewGetNotificationsByUserUseCase(notificationRepository)
	updateNotificationUseCase := usecase.NewUpdateNotificationUseCase(notificationRepository, userRepository)
	deleteNotificationUseCase := usecase.NewDeleteNotification(notificationRepository)

	dependencies := routes.AppDependencies{
		Router:                       router,
		Cfg:                          cfg,
		CreateUserUseCase:            createUserUseCase,
		GetUserUseCase:               getUserUseCase,
		UpdateUserUseCase:            updateUserUseCase,
		CreateNotificationUseCase:    createNotificationUseCase,
		GetNotificationByUserUseCase: getNotificationsByUserUseCase,
		UpdateNotificationUseCase:    updateNotificationUseCase,
		DeleteNotificationUseCase:    deleteNotificationUseCase,
	}

	routes.SetupRoutes(dependencies)

	log.Fatal(router.Run())
}
