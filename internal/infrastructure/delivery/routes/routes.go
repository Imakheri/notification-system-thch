package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers"
	"github.com/imakheri/notifications-thch/internal/infrastructure/middleware"
)

type AppDependencies struct {
	Router                       *gin.Engine
	Cfg                          *config.Config
	CreateUserUseCase            usecase.CreateUserUseCase
	GetUserUseCase               usecase.GetUserUseCase
	UpdateUserUseCase            usecase.UpdateUserUseCase
	CreateNotificationUseCase    usecase.CreateNotificationUseCase
	GetNotificationByUserUseCase usecase.GetNotificationsByUserUseCase
	UpdateNotificationUseCase    usecase.UpdateNotificationUseCase
	DeleteNotificationUseCase    usecase.DeleteNotificationUseCase
}

func SetupRoutes(dps AppDependencies) {
	v1 := dps.Router.Group("/api/v1")
	{
		public := v1.Group("/")
		public.Use(middleware.NewAPIKey(dps.Cfg).ValidateAPIKey())
		{
			public.POST("/users", handlers.CreateUserHandler(dps.CreateUserUseCase))
			public.POST("/user", handlers.GetUserHandler(dps.GetUserUseCase))
		}

		auth := v1.Group("/")
		auth.Use(
			middleware.NewAPIKey(dps.Cfg).ValidateAPIKey(),
			middleware.NewAuthorizeJWT(dps.Cfg).AuthorizeJWT(),
		)
		{
			auth.PUT("/user", handlers.UpdateUserHandler(dps.UpdateUserUseCase))
			auth.POST("/notification", handlers.CreateNotificationHandler(dps.CreateNotificationUseCase))
			auth.GET("/notifications", handlers.GetNotificationsByUserIDHandler(dps.GetNotificationByUserUseCase))
			auth.PUT("/notification/:id", handlers.UpdateNotificationHandler(dps.UpdateNotificationUseCase))
			auth.DELETE("/notification/:id", handlers.DeleteNotificationHandler(dps.DeleteNotificationUseCase))
		}
	}
}
