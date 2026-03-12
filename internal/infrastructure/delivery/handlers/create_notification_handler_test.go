package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateNotificationHandler(t *testing.T) {
	var fixedTime = time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	type fields struct {
		createNotificationUseCase func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase
	}
	type args struct {
		userID         uint
		userEmail      string
		hasAuthContext bool
		body           string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody string
	}{
		{
			name: "Create Notification Handler",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(2), "williams@example.com", entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								Email: "charlesd@example.com",
							},
						},
					}).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    &fixedTime,
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								ID:            3,
								Name:          "Charles Dickens",
								Email:         "charlesd@example.com",
								Phone:         "0987654321",
								DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
								Notifications: nil,
							},
						},
					}, nil)
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "williams@example.com",
				hasAuthContext: true,
				body: `{
						"title": "Test Notification",
						"content": "This is a test notification",
						"channel_id": 1,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusCreated,
			wantBody: `{
						"notification": {
							"id": 1,
							"created_by": 2,
							"sent_at": "2026-03-09T12:00:00Z",		
							"title": "Test Notification",
							"content": "This is a test notification",
							"channel_id": 1,
							"recipients": [
								{
									"id": 3,
									"email": "charlesd@example.com"
								}
							]
						}
					}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "williams@example.com",
				hasAuthContext: false,
				body: `{
						"title": "Test Notification",
						"content": "This is a test notification",
						"channel_id": ,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"invalid character ',' looking for beginning of value"}`,
		},
		{
			name: "Error getting user email from token",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "",
				hasAuthContext: true,
				body: `{
						"title": "Test Notification",
						"content": "This is a test notification",
						"channel_id": 1,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error": "can not get user email"}`,
		},
		{
			name: "No content on notification error",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "williams@example.com",
				hasAuthContext: true,
				body: `{
						"title": "Test Notification",
						"content": "",
						"channel_id": 1,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"notification must have a content"}`,
		},
		{
			name: "No content on notification error",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "williams@example.com",
				hasAuthContext: true,
				body: `{
						"title": "Test Notification",
						"content": "",
						"channel_id": 1,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"notification must have a content"}`,
		},
		{
			name: "Create Notification Handler",
			fields: fields{
				createNotificationUseCase: func(ctrl *gomock.Controller) usecase.CreateNotificationUseCase {
					m := mocks.NewMockCreateNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(2), "williams@example.com", entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								Email: "charlesd@example.com",
							},
						},
					}).Return(entities.Notification{}, errors.New("db error: connection lost"))
					return m
				},
			},
			args: args{
				userID:         2,
				userEmail:      "williams@example.com",
				hasAuthContext: true,
				body: `{
						"title": "Test Notification",
						"content": "This is a test notification",
						"channel_id": 1,
						"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}`,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error": "db error: connection lost"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tt.args.hasAuthContext {
				ctx.Set("id", tt.args.userID)
				ctx.Set("email", tt.args.userEmail)
			}

			ctx.Request, _ = http.NewRequest(http.MethodPost, "/notification", strings.NewReader(tt.args.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockUseCase := tt.fields.createNotificationUseCase(ctrl)
			handler := CreateNotificationHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
