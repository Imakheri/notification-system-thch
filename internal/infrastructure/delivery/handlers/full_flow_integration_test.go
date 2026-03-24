package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure"
	handler_dtos "github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestE2E_FullFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := TestDbStruct.DatabaseConnection.Begin()
	defer tx.Rollback()

	dbWrapper := &repository.Database{DatabaseConnection: tx}
	userRepository := repository.NewUserRepository(dbWrapper)
	notificationRepository := repository.NewNotificationRepository(dbWrapper)
	channelRepository := repository.NewChannelRepository(dbWrapper)
	mockSimulatedApiService := mocks.NewMockSimulatedApiService(ctrl)
	clock := infrastructure.RealClock{}

	createUserUseCase := usecase.NewCreateUserUseCase(userRepository)
	createNotificationUseCase := usecase.NewCreateNotificationUseCase(notificationRepository, userRepository, channelRepository, mockSimulatedApiService, clock)
	getNotificationByUserID := usecase.NewGetNotificationsByUserUseCase(notificationRepository)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	userInput := handler_dtos.CreateUserDTO{
		Name:        "Charles Dickens",
		Password:    "Christmas123.",
		Email:       "charlesd@example.com",
		Phone:       "0987654321",
		DeviceToken: "e0d8c7b6a5y4t3r2e1q0",
	}

	body, _ := json.Marshal(userInput)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	ctx.Request.Header.Add("Content-Type", "application/json")

	CreateUserHandler(createUserUseCase)(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)

	userInput = handler_dtos.CreateUserDTO{
		Name:        "William Shakespeare",
		Password:    "Hamlet123.",
		Email:       "williams@example.com",
		Phone:       "1234567890",
		DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
	}

	body, _ = json.Marshal(userInput)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	ctx.Request.Header.Add("Content-Type", "application/json")

	CreateUserHandler(createUserUseCase)(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)

	var userInDB dtos.UserModel
	err := tx.Where("email = ?", "williams@example.com").First(&userInDB).Error

	assert.NoError(t, err, "User should have been created on database")
	assert.NotZero(t, userInDB.ID, "User ID should not be zero")

	createdUserID := userInDB.ID
	createdUserEmail := userInDB.Email

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Set("user_id", createdUserID)
	ctx.Set("email", createdUserEmail)
	notificationInput := handler_dtos.CreateNotificationDTO{
		Title:     "Test Notification",
		Content:   "This is a test notification",
		ChannelID: 1,
		Recipients: []handler_dtos.UserIntoNotificationDTO{
			{
				Email: "charlesd@example.com",
			},
		},
	}

	bodyNotification, _ := json.Marshal(notificationInput)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/notification", bytes.NewBuffer(bodyNotification))
	ctx.Request.Header.Add("Content-Type", "application/json")

	mockSimulatedApiService.EXPECT().RandomizeHTTPStatus().Return(http.StatusOK, nil)

	CreateNotificationHandler(createNotificationUseCase)(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Set("user_id", createdUserID)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/notifications", nil)

	GetNotificationsByUserIDHandler(getNotificationByUserID)(ctx)
	assert.Equal(t, http.StatusOK, w.Code)

	var notificationResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notificationResponse)
	notifications := notificationResponse["notifications"].([]interface{})

	assert.NotEmpty(t, notifications)
	firstNotification := notifications[0].(map[string]interface{})
	assert.Equal(t, "Test Notification", firstNotification["title"])
	assert.Equal(t, "This is a test notification", firstNotification["content"])
	assert.Equal(t, 1, int(firstNotification["channel_id"].(float64)))
}
