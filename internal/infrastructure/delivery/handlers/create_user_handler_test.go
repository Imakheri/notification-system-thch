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

func TestCreateUserHandler(t *testing.T) {
	type fields struct {
		createUserUseCase func(ctrl *gomock.Controller) usecase.CreateUserUseCase
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
			name: "Created user successfully",
			fields: fields{
				createUserUseCase: func(ctrl *gomock.Controller) usecase.CreateUserUseCase {
					m := mocks.NewMockCreateUserUseCase(ctrl)
					m.EXPECT().Exec(entities.User{
						Name:        "William Shakespeare",
						Password:    "Hamlet123.",
						Email:       "williams@example.com",
						Phone:       "1234567890",
						DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
					}).Return(entities.User{
						ID:          1,
						Name:        "William Shakespeare",
						Password:    "Hamlet123.",
						Email:       "williams@example.com",
						Phone:       "1234567890",
						DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
					}, nil)
					return m
				},
			},
			args: args{
				body: `{
  					"name": "William Shakespeare",
  					"email": "williams@example.com",
					"password": "Hamlet123.",
  					"phone": "1234567890",
  					"device_token": "0q1e2r3t4y5a6b7c8d9e"
				}`,
			},
			wantStatus: http.StatusCreated,
			wantBody: `{
							"user": {
								"name": "William Shakespeare",
								"email": "williams@example.com",
								"phone": "1234567890",
								"device_token": "0q1e2r3t4y5a6b7c8d9e"
							}
						}`,
		},
		{
			name: "Error triggered by malformed body (JSON)",
			fields: fields{
				createUserUseCase: func(ctrl *gomock.Controller) usecase.CreateUserUseCase {
					m := mocks.NewMockCreateUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				body: `{
  					"name": "William Shakespeare",
  					"email": "williams@example.com",
					"password": "ACommaIsMissingHere!123"
  					"phone": "1234567890",
  					"device_token": "0q1e2r3t4y5a6b7c8d9e"
				}`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid character '\"' after object key:value pair"}`,
		},
		{
			name: "Email invalid structure",
			fields: fields{
				createUserUseCase: func(ctrl *gomock.Controller) usecase.CreateUserUseCase {
					m := mocks.NewMockCreateUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				body: `{
  					"name": "William Shakespeare",
  					"email": "invalidemailexamplecom",
					"password": "Hamlet123.",
  					"phone": "1234567890",
  					"device_token": "0q1e2r3t4y5a6b7c8d9e"
				}`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid email structure"}`,
		},
		{
			name: "Insecure password entered",
			fields: fields{
				createUserUseCase: func(ctrl *gomock.Controller) usecase.CreateUserUseCase {
					m := mocks.NewMockCreateUserUseCase(ctrl)
					return m
				},
			},
			args: args{
				body: `{
  					"name": "William Shakespeare",
  					"email": "invalidemailexamplecom",
					"password": "password",
  					"phone": "1234567890",
  					"device_token": "0q1e2r3t4y5a6b7c8d9e"
				}`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"password must contain at least one upper case and one lower case character, one digit and one special character"}`,
		},
		{
			name: "User already exists",
			fields: fields{
				createUserUseCase: func(ctrl *gomock.Controller) usecase.CreateUserUseCase {
					m := mocks.NewMockCreateUserUseCase(ctrl)
					m.EXPECT().Exec(entities.User{
						Name:        "William Shakespeare",
						Password:    "Hamlet123.",
						Email:       "williams@example.com",
						Phone:       "1234567890",
						DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
					}).Return(entities.User{}, errors.New("user already exists"))
					return m
				},
			},
			args: args{
				body: `{
  					"name": "William Shakespeare",
  					"email": "williams@example.com",
					"password": "Hamlet123.",
  					"phone": "1234567890",
  					"device_token": "0q1e2r3t4y5a6b7c8d9e"
				}`,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"user already exists"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.args.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockUseCase := tt.fields.createUserUseCase(ctrl)
			handler := CreateUserHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
