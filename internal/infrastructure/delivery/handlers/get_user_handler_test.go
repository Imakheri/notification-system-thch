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

func TestGetUserHandler(t *testing.T) {
	type fields struct {
		getUserUseCase func(ctrl *gomock.Controller) usecase.GetUserUseCase
	}
	type args struct {
		body string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "Got user successfully (Login successful) ",
			fields: fields{
				getUserUseCase: func(ctrl *gomock.Controller) usecase.GetUserUseCase {
					m := mocks.NewMockGetUserUseCase(ctrl)
					m.EXPECT().Exec(entities.User{
						Email:    "williams@example.com",
						Password: "Hamlet123.",
					}).Return(entities.User{
						ID:            1,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, "", nil)
					return m
				},
			},
			args: args{
				body: `{"email": "williams@example.com", "password": "Hamlet123."}`,
			},
			wantStatus: http.StatusOK,
			wantBody: `{
						  "token": "",
						  "user": {
							"name": "William Shakespeare",
							"email": "williams@example.com",
							"phone": "1234567890",
							"device_token": "0q1e2r3t4y5a6b7c8d9e"
						  }
						}`,
		},
		{
			name: "User does not exists (Login unsuccessful) ",
			fields: fields{
				getUserUseCase: func(ctrl *gomock.Controller) usecase.GetUserUseCase {
					m := mocks.NewMockGetUserUseCase(ctrl)
					m.EXPECT().Exec(entities.User{
						Email:    "example@example.com",
						Password: "ThisIsAValidPassword123.",
					}).Return(entities.User{}, "", errors.New("the e-mail address or password is incorrect"))
					return m
				},
			},
			args: args{
				body: `{"email": "example@example.com", "password": "ThisIsAValidPassword123."}`,
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"error": "the e-mail address or password is incorrect"}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				getUserUseCase: func(ctrl *gomock.Controller) usecase.GetUserUseCase {
					m := mocks.NewMockGetUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				body: `{"email": williams@example.com", "password": "Hamlet123."}`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid character 'w' looking for beginning of value"}`,
		},
		{
			name: "UseCase failed because database error",
			fields: fields{
				getUserUseCase: func(ctrl *gomock.Controller) usecase.GetUserUseCase {
					m := mocks.NewMockGetUserUseCase(ctrl)
					m.EXPECT().Exec(entities.User{
						Email:    "williams@example.com",
						Password: "Hamlet123.",
					}).Return(entities.User{
						ID:            1,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, "", errors.New("sql: Scan error on column 'content': incompatible type"))
					return m
				},
			},
			args: args{
				body: `{"email": "williams@example.com", "password": "Hamlet123."}`,
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"error": "sql: Scan error on column 'content': incompatible type"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request, _ = http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(tt.args.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockUseCase := tt.fields.getUserUseCase(ctrl)
			handler := GetUserHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
