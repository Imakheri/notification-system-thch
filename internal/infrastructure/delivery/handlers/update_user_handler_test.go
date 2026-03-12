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

func TestUpdateUserHandler(t *testing.T) {
	type fields struct {
		updateUserUseCase func(ctrl *gomock.Controller) usecase.UpdateUserUseCase
	}
	type args struct {
		authUserEmail  string
		hasAuthContext bool
		body           string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "User updated successfully",
			fields: fields{
				updateUserUseCase: func(ctrl *gomock.Controller) usecase.UpdateUserUseCase {
					m := mocks.NewMockUpdateUserUseCase(ctrl)
					m.EXPECT().Exec("fernandop@example.com", entities.User{
						Name:        "Ricardo Reis",
						DeviceToken: "12345qwerty67890",
					}).Return(entities.User{
						ID:          3,
						Name:        "Ricardo Reis",
						Email:       "fernandop@example.com",
						Phone:       "0129834765",
						DeviceToken: "12345qwerty67890",
					}, nil)
					return m
				},
			},
			args: args{
				authUserEmail:  "fernandop@example.com",
				hasAuthContext: true,
				body: `{
						"name": "Ricardo Reis",
						"device_token": "12345qwerty67890"
						}`,
			},
			wantStatus: http.StatusCreated,
			wantBody: `{
							"user": {
								"name": "Ricardo Reis",
								"email": "fernandop@example.com",
 								"phone": "0129834765",
								"device_token": "12345qwerty67890"
							}
						}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				updateUserUseCase: func(ctrl *gomock.Controller) usecase.UpdateUserUseCase {
					m := mocks.NewMockUpdateUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				authUserEmail:  "fernandop@example.com",
				hasAuthContext: false,
				body: `{
						"name" "Ricardo Reis",
						"device_token": "12345qwerty67890"
						}`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error": "invalid character '\"' after object key"}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				updateUserUseCase: func(ctrl *gomock.Controller) usecase.UpdateUserUseCase {
					m := mocks.NewMockUpdateUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				authUserEmail:  "",
				hasAuthContext: false,
				body: `{
						"name": "Ricardo Reis",
						"device_token": "12345qwerty67890"
						}`,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error": "can not get user email"}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				updateUserUseCase: func(ctrl *gomock.Controller) usecase.UpdateUserUseCase {
					m := mocks.NewMockUpdateUserUseCase(ctrl)
					m.EXPECT().Exec("fernandop@example.com", entities.User{
						Name:        "Ricardo Reis",
						DeviceToken: "12345qwerty67890",
					}).Return(entities.User{}, errors.New("db error: connection lost during update"))
					return m
				},
			},
			args: args{
				authUserEmail:  "fernandop@example.com",
				hasAuthContext: true,
				body: `{
						"name": "Ricardo Reis",
						"device_token": "12345qwerty67890"
						}`,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error": "db error: connection lost during update"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tt.args.hasAuthContext {
				ctx.Set("email", tt.args.authUserEmail)
			}

			ctx.Request, _ = http.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(tt.args.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockUseCase := tt.fields.updateUserUseCase(ctrl)
			handler := UpdateUserHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())

		})
	}
}
