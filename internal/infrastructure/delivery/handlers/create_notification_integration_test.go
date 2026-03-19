//go:build integration

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
	handlers_dtos "github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestIntegrationCreateNotification(t *testing.T) {
	type fields struct {
		setupData func(tx *gorm.DB) (uint, uint)
	}
	type args struct {
		input gin.H
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		userID         interface{}
		userEmail      interface{}
		expectedStatus int
	}{
		{
			name: "Success_Create_Notification",
			fields: fields{
				setupData: func(tx *gorm.DB) (uint, uint) {
					user := dtos.UserModel{
						Name:        "William Shakespeare",
						Email:       "williams@example.com",
						Phone:       "1234567890",
						DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
					}
					tx.Create(&user)
					recipient := dtos.UserModel{
						Name:        "Charles Dickens",
						Email:       "charlesd@example.com",
						Phone:       "0987654321",
						DeviceToken: "e0d8c7b6a5y4t3r2e1q0",
					}
					tx.Create(&recipient)
					return user.ID, recipient.ID
				},
			},
			args: args{
				input: gin.H{
					"title":      "Test Notification",
					"content":    "This is a test notification",
					"channel_id": 1,
					"recipients": []dtos.UserModel{
						{
							Email: "charlesd@example.com",
						},
					},
				},
			},
			userID:         "",
			userEmail:      "williams@example.com",
			expectedStatus: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSimulatedApiService := mocks.NewMockSimulatedApiService(ctrl)
			mockSimulatedApiService.EXPECT().RandomizeHTTPStatus().Return(http.StatusOK, nil)

			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			var userID, _ uint
			if tt.fields.setupData != nil {
				userID, _ = tt.fields.setupData(tx)
			}

			dbWrapper := &repository.Database{DatabaseConnection: tx}
			notificationRepository := repository.NewNotificationRepository(dbWrapper)
			userRepository := repository.NewUserRepository(dbWrapper)
			channelRepository := repository.NewChannelRepository(dbWrapper)
			clock := infrastructure.RealClock{}
			useCase := usecase.NewCreateNotificationUseCase(notificationRepository, userRepository, channelRepository, mockSimulatedApiService, clock)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Set("id", userID)
			ctx.Set("email", tt.userEmail)

			jsonBody, _ := json.Marshal(tt.args.input)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/notification", bytes.NewBuffer(jsonBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := CreateNotificationHandler(useCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var notificationsInDB []dtos.NotificationModel
				tx.Preload("Recipients").Where("created_by = ?", userID).Find(&notificationsInDB)

				var actualResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &actualResponse)
				actualList := actualResponse["notifications"].([]interface{})

				var expectedList []interface{}
				for _, n := range notificationsInDB {
					dto := handlers_dtos.NotificationToDto(dtos.NotificationModelToEntity(n))
					dtoJSON, _ := json.Marshal(dto)
					var dtoMap map[string]interface{}
					json.Unmarshal(dtoJSON, &dtoMap)
					delete(dtoMap, "sent_at")
					expectedList = append(expectedList, dtoMap)
				}

				for _, a := range actualList {
					delete(a.(map[string]interface{}), "sent_at")
				}

				assert.Equal(t, expectedList, actualList, "Notification list does not match")
			}
		})
	}
}
