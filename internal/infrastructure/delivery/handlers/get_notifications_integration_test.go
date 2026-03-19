//go:build integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	handlers_dtos "github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegrationGetNotificationsByUserID(t *testing.T) {
	type fields struct {
		setupData func(tx *gorm.DB) (uint, uint)
	}
	tests := []struct {
		name           string
		fields         fields
		expectedStatus int
	}{
		{
			name: "Success_Get_Notifications",
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
					notifications := []dtos.NotificationModel{
						{
							CreatedBy: user.ID,
							SentAt:    nil,
							Title:     "Test Notification",
							Content:   "This is a test notification",
							ChannelID: 1,
							Recipients: []dtos.UserModel{
								{
									Email: "charlesd@example.com",
								},
							},
						},
						{
							CreatedBy: user.ID,
							SentAt:    nil,
							Title:     "Test Notification 2",
							Content:   "This is a test notification 2",
							ChannelID: 2,
							Recipients: []dtos.UserModel{
								{
									Email: "charlesd@example.com",
								},
							},
						},
					}
					for _, n := range notifications {
						tx.Create(&n)
					}
					return user.ID, recipient.ID
				},
			},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			var userID, _ uint
			if tt.fields.setupData != nil {
				userID, _ = tt.fields.setupData(tx)
			}

			dbWrapper := &repository.Database{DatabaseConnection: tx}
			notificationRepository := repository.NewNotificationRepository(dbWrapper)
			getNotificationByUserUseCase := usecase.NewGetNotificationsByUserUseCase(notificationRepository)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/notifications", nil)
			ctx.Set("id", userID)

			h := GetNotificationsByUserIDHandler(getNotificationByUserUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var notificationsInDB []dtos.NotificationModel
				tx.Preload("Recipients").Where("created_by = ?", userID).Find(&notificationsInDB)

				var expectedDTOs []handlers_dtos.NotificationResponseDTO
				for _, n := range notificationsInDB {
					expectedDTOs = append(expectedDTOs, handlers_dtos.NotificationToDto(dtos.NotificationModelToEntity(n)))
				}

				expectedResponse := gin.H{
					"notifications": expectedDTOs,
				}

				expectedJSON, _ := json.Marshal(expectedResponse)

				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			}
		})
	}
}
