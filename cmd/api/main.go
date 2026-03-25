package main

import (
	_ "github.com/imakheri/notifications-thch/cmd/api/docs"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/routes"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/service"
)

// @title           Notification Manager System - THC
// @version         0.1
// @description     Service for managing and sending notifications (Take Home Challenge)
// @termsOfService http://swagger.io/terms/

// @contact.name   Julián Rentería
// @contact.email   imakheri@gmail.com

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api_key

func main() {
	router := gin.Default()
	cfg := config.Load()
	db := repository.NewDatabase(cfg)
	realClock := infrastructure.RealClock{}

	userRepository := repository.NewUserRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	channelRepository := repository.NewChannelRepository(db)

	simulatedApiService := service.NewSimulatedApiService()
	jwtService := service.NewJWTService(cfg)

	createUserUseCase := usecase.NewCreateUserUseCase(userRepository)
	getUserUseCase := usecase.NewGetUserUseCase(userRepository, jwtService)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepository)

	createNotificationUseCase := usecase.NewCreateNotificationUseCase(notificationRepository, userRepository, channelRepository, simulatedApiService, realClock)
	getNotificationsByUserUseCase := usecase.NewGetNotificationsByUserUseCase(notificationRepository)
	updateNotificationUseCase := usecase.NewUpdateNotificationUseCase(notificationRepository, userRepository)
	deleteNotificationUseCase := usecase.NewDeleteNotificationUseCase(notificationRepository)

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
