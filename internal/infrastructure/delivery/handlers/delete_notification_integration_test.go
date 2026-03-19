//go:build integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegrationDeleteNotification(t *testing.T) {
	type args struct {
		notificationID string
		userID         uint
		hasAuthContext bool
	}
	tests := []struct {
		name           string
		setupData      func(tx *gorm.DB) string
		args           args
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Notification deleted successfully (IT)",
			setupData: func(tx *gorm.DB) string {
				notification := dtos.NotificationModel{
					CreatedBy: 1,
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
				return strconv.Itoa(int(notification.ID))
			},
			args: args{
				userID:         1,
				hasAuthContext: true,
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"message": "Notification deleted successfully",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			if tt.setupData != nil {
				tt.args.notificationID = tt.setupData(tx)
			}

			dbWrapper := &repository.Database{
				DatabaseConnection: tx,
			}

			notificationRepository := repository.NewNotificationRepository(dbWrapper)
			deleteNotificationUseCase := usecase.NewDeleteNotificationUseCase(notificationRepository)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{{Key: "id", Value: tt.args.notificationID}}
			if tt.args.hasAuthContext {
				ctx.Set("id", tt.args.userID)
			}
			ctx.Request, _ = http.NewRequest(http.MethodDelete, "/notification/:id", nil)

			h := DeleteNotificationHandler(deleteNotificationUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var countNormal int64
				var countPhysical int64
				var deletedNotification dtos.NotificationModel

				tx.Model(&dtos.NotificationModel{}).Where("id = ?", tt.args.notificationID).Count(&countNormal)
				assert.Equal(t, int64(0), countNormal, "notification should not be visible in normal queries")

				tx.Unscoped().Model(&dtos.NotificationModel{}).Where("id = ?", tt.args.notificationID).Count(&countPhysical)
				assert.Equal(t, int64(1), countPhysical, "physical record should exists (Soft Delete)")

				tx.Unscoped().First(&deletedNotification, tt.args.notificationID)
				assert.NotNil(t, deletedNotification.DeletedAt, "DeletedAt column should have information")
			}

			var response map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			for key, val := range tt.expectedBody {
				assert.Equal(t, val, response[key], "%s does not match", key)
			}
		})
	}
}
