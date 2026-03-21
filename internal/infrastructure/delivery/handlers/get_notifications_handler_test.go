package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetNotificationsByUserIDHandler(t *testing.T) {
	var fixedTime = time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	type fields struct {
		GetNotificationByUserUseCase func(ctrl *gomock.Controller) usecase.GetNotificationsByUserUseCase
	}
	type args struct {
		authUserID     uint
		hasAuthContext bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "Got all notifications from user successfully",
			fields: fields{
				GetNotificationByUserUseCase: func(ctrl *gomock.Controller) usecase.GetNotificationsByUserUseCase {
					m := mocks.NewMockGetNotificationsByUserUseCase(ctrl)
					m.EXPECT().Exec(uint(1)).Return([]entities.Notification{
						{
							ID:        1,
							CreatedBy: 1,
							SentAt:    &fixedTime,
							Title:     "Test Notification",
							Content:   "This is a test notification",
							ChannelID: 2,
							Recipients: []entities.User{
								{
									ID:            2,
									Name:          "William Shakespeare",
									Email:         "williams@example.com",
									Phone:         "1234567890",
									DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
									Notifications: nil,
								},
							},
						},
						{
							ID:        2,
							CreatedBy: 1,
							SentAt:    &fixedTime,
							Title:     "Test Notification #2",
							Content:   "This is a test notification number 2",
							ChannelID: 3,
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
						},
					}, nil)
					return m
				},
			},
			args: args{
				authUserID:     1,
				hasAuthContext: true,
			},
			wantStatus: http.StatusOK,
			wantBody: `{
						  "notifications": [
							{
							  "id": 1,
							  "created_by": 1,
							  "sent_at": "2026-03-09T12:00:00Z",
							  "title": "Test Notification",
							  "content": "This is a test notification",
							  "channel_id": 2,
							  "recipients": [
								{
								  "id": 2,
								  "email": "williams@example.com"
								}
							  ]
							},
							{
							  "id": 2,
							  "created_by": 1,
							  "sent_at": "2026-03-09T12:00:00Z",
							  "title": "Test Notification #2",
							  "content": "This is a test notification number 2",
							  "channel_id": 3,
							  "recipients": [
								{
								  "id": 3,
								  "email": "charlesd@example.com"
								}
							  ]
							}
						  ]
						}`,
		},
		{
			name: "User has no notification, but response successfully",
			fields: fields{
				GetNotificationByUserUseCase: func(ctrl *gomock.Controller) usecase.GetNotificationsByUserUseCase {
					m := mocks.NewMockGetNotificationsByUserUseCase(ctrl)
					m.EXPECT().Exec(uint(1)).Return([]entities.Notification{}, nil)
					return m
				},
			},
			args: args{
				authUserID:     1,
				hasAuthContext: true,
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"notifications": []}`,
		},
		{
			name: "Error getting user id from token",
			fields: fields{
				GetNotificationByUserUseCase: func(ctrl *gomock.Controller) usecase.GetNotificationsByUserUseCase {
					m := mocks.NewMockGetNotificationsByUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				authUserID:     0,
				hasAuthContext: false,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"can not get user id"}`,
		},
		{
			name: "UseCase failed because database error ",
			fields: fields{
				GetNotificationByUserUseCase: func(ctrl *gomock.Controller) usecase.GetNotificationsByUserUseCase {
					m := mocks.NewMockGetNotificationsByUserUseCase(ctrl)
					m.EXPECT().Exec(uint(1)).Return([]entities.Notification{}, errors.New("sql: Scan error on column 'content': incompatible type"))
					return m
				},
			},
			args: args{
				authUserID:     1,
				hasAuthContext: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error": "sql: Scan error on column 'content': incompatible type"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tt.args.hasAuthContext {
				ctx.Set("user_id", tt.args.authUserID)
			}

			mockUseCase := tt.fields.GetNotificationByUserUseCase(ctrl)
			handler := GetNotificationsByUserIDHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
