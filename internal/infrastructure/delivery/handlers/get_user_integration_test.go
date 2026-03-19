//go:build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/delivery/handlers/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	repository_dtos "github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/imakheri/notifications-thch/internal/infrastructure/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegrationGetUser(t *testing.T) {
	type fields struct {
		setupData func(tx *gorm.DB)
	}
	type args struct {
		input dtos.GetUserDTO
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedStatus   int
		expectedUserBody string
	}{
		{
			name: "User successfully logged in",
			fields: fields{
				setupData: func(tx *gorm.DB) {
					user := repository_dtos.UserModel{
						Name:        "William Shakespeare",
						Email:       "williams@example.com",
						Password:    "$2a$10$eQMJB/hVz5wpVN/nNmEyIenP/gv8B09TlNSXkBue7j4UlgUfgrxK.",
						Phone:       "1234567890",
						DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
					}
					tx.Create(&user)
				},
			},
			args: args{
				input: dtos.GetUserDTO{
					Email:    "williams@example.com",
					Password: "Hamlet123.",
				},
			},
			expectedStatus: http.StatusOK,
			expectedUserBody: `{
				"name":  "William Shakespeare",
				"email": "williams@example.com",
				"phone": "1234567890",	
				"device_token": "0q1e2r3t4y5a6b7c8d9e"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			if tt.fields.setupData != nil {
				tt.fields.setupData(tx)
			}

			dbWrapper := &repository.Database{DatabaseConnection: tx}
			jwtService := service.NewJWTService(&config.Config{})
			userRepository := repository.NewUserRepository(dbWrapper)
			getUserUseCase := usecase.NewGetUserUseCase(userRepository, jwtService)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			jsonInput, _ := json.Marshal(tt.args.input)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonInput))
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := GetUserHandler(getUserUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.NotNil(t, response["token"])
				assert.NotEmpty(t, response["token"])

				userJson, _ := json.Marshal(response["user"])
				assert.JSONEq(t, tt.expectedUserBody, string(userJson))
			}
		})
	}
}
