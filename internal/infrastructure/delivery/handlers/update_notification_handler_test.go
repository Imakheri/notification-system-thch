package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateNotificationHandler(t *testing.T) {
	type fields struct {
		updateNotificationUseCase func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase
	}
	type args struct {
		notificationID string
		userID         uint
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
			name: "Update notification successfully",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(2), 1, entities.Notification{
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly.",
						ChannelID: 3,
					}).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly.",
						ChannelID: 3,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
						},
					}, nil)
					return m
				},
			},
			args: args{
				notificationID: "1",
				userID:         2,
				hasAuthContext: true,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 3
						}`,
			},
			wantCode: http.StatusOK,
			wantBody: `{
						"notification": {
							"id": 1,
							"created_by": 2,
							"sent_at": null,		
							"title": "Earthquake detected in your area!",
							"content": "Drop, cover, and hold on. Please proceed calmly.",
							"channel_id": 3,
							"recipients": [
								{
									"email": "charlesd@example.com"
								}
							]
						}
					}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "1",
				userID:         0,
				hasAuthContext: false,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id: 3
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error": "invalid character '\\n' in string literal"}`,
		},
		{
			name: "Missing notification ID from param",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "",
				userID:         0,
				hasAuthContext: false,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 3
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error": "must enter a valid id"}`,
		},
		{
			name: "Invalid notification ID param",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "qwerty",
				userID:         0,
				hasAuthContext: false,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 3
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error": "strconv.Atoi: parsing \"qwerty\": invalid syntax"}`,
		},
		{
			name: "Error getting id from token",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "1",
				userID:         0,
				hasAuthContext: false,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 3
						}`,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error": "can not get user id"}`,
		},
		{
			name: "Invalid channel entered",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "1",
				userID:         2,
				hasAuthContext: true,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 5
						}`,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error": "must use a valid channel"}`,
		},
		{
			name: "Update notification successfully",
			fields: fields{
				updateNotificationUseCase: func(ctrl *gomock.Controller) usecase.UpdateNotificationUseCase {
					m := mocks.NewMockUpdateNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(2), 1, entities.Notification{
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly.",
						ChannelID: 3,
					}).Return(entities.Notification{}, errors.New("db: no rows in result set"))
					return m
				},
			},
			args: args{
				notificationID: "1",
				userID:         2,
				hasAuthContext: true,
				body: `{
						"title": "Earthquake detected in your area!",
						"content": "Drop, cover, and hold on. Please proceed calmly.",
						"channel_id": 3
						}`,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error": "db: no rows in result set"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = gin.Params{{Key: "id", Value: tt.args.notificationID}}

			if tt.args.hasAuthContext {
				ctx.Set("user_id", tt.args.userID)
			}

			ctx.Request, _ = http.NewRequest(http.MethodPut, "/notification/:id", bytes.NewBufferString(tt.args.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockUseCase := tt.fields.updateNotificationUseCase(ctrl)
			handler := UpdateNotificationHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())

		})
	}
}
