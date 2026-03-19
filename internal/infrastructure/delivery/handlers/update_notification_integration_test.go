//go:build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	handlers_dtos "github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegrationUpdateNotification(t *testing.T) {
	type fields struct {
		setupData func(tx *gorm.DB) (uint, uint, uint)
	}
	type args struct {
		notificationID string
		input          gin.H
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "Notification updated successfully",
			fields: fields{
				setupData: func(tx *gorm.DB) (uint, uint, uint) {
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
					notification := dtos.NotificationModel{
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
					}
					tx.Create(&notification)
					return user.ID, recipient.ID, notification.ID
				},
			},
			args: args{
				notificationID: "1",
				input: gin.H{
					"title":      "Valid notification title",
					"content":    "This is a valid notification content",
					"channel_id": 2,
					"recipients": []dtos.UserModel{
						{
							Email: "charlesd@example.com",
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: `{
								"notification": {
									"title": Valid notification title",
									"content": "This is a valid notification content",
									"created_by": 1,
									"sent_at": null,
									"channel_id": 2,
									"recipients": [
										{
											"id": 
											"email": "charlesd@example.com"
										}
									]
								}
							}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			var userID, _, notificationID uint
			if tt.fields.setupData != nil {
				userID, _, notificationID = tt.fields.setupData(tx)
			}

			dbWrapper := &repository.Database{DatabaseConnection: tx}
			notificationRepository := repository.NewNotificationRepository(dbWrapper)
			userRepository := repository.NewUserRepository(dbWrapper)
			updateNotificationUseCase := usecase.NewUpdateNotificationUseCase(notificationRepository, userRepository)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			paramID := tt.args.notificationID
			if notificationID > 0 {
				paramID = strconv.FormatUint(uint64(notificationID), 10)
			}

			ctx.Params = []gin.Param{{Key: "id", Value: paramID}}

			ctx.Set("id", userID)

			jsonInput, _ := json.Marshal(tt.args.input)
			ctx.Request, _ = http.NewRequest(http.MethodPut, "/notifications/"+paramID, bytes.NewBuffer(jsonInput))
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := UpdateNotificationHandler(updateNotificationUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var updatedInDB dtos.NotificationModel
				tx.First(&updatedInDB, notificationID)

				expectedDto := dtos.NotificationModelToEntity(updatedInDB)
				expectedDto2 := handlers_dtos.NotificationToDto(expectedDto)

				expectedResponse := gin.H{
					"notification": expectedDto2,
				}

				expectedJSON, _ := json.Marshal(expectedResponse)
				assert.JSONEq(t, string(expectedJSON), w.Body.String(), "JSON response is wrong")
			}

			var updatedNotificationInDB dtos.NotificationModel
			tx.First(&updatedNotificationInDB, notificationID)
			if newTitle, ok := tt.args.input["title"].(string); ok {
				assert.Equal(t, newTitle, updatedNotificationInDB.Title, "Notification title was not updated on database")
			}
		})
	}
}
