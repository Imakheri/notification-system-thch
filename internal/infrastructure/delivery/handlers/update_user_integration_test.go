//go:build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegrationUpdateUser(t *testing.T) {
	type fields struct {
		setupData func(tx *gorm.DB)
	}
	type args struct {
		userEmail      string
		input          gin.H
		hasAuthContext bool
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "User successfully updated",
			fields: fields{
				setupData: func(tx *gorm.DB) {
					user := dtos.UserModel{
						Name:        "Fernando Pessoa",
						Password:    "Mensagem123.",
						Email:       "fernandop@example.com",
						Phone:       "0129834765",
						DeviceToken: "",
					}
					tx.Create(&user)
				},
			},
			args: args{
				userEmail: "fernandop@example.com",
				input: gin.H{
					"name":         "Ricardo Reis",
					"device_token": "12345qwerty67890",
				},
				hasAuthContext: true,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
							  "user": {
								"name":  "Ricardo Reis",
								"email": "fernandop@example.com",
								"phone": "0129834765",
								"device_token": "12345qwerty67890"
							  }
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
			userRepository := repository.NewUserRepository(dbWrapper)
			updateUserUseCase := usecase.NewUpdateUserUseCase(userRepository)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tt.args.hasAuthContext {
				ctx.Set("email", tt.args.userEmail)
			}

			jsonInput, _ := json.Marshal(tt.args.input)
			ctx.Request, _ = http.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(jsonInput))
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := UpdateUserHandler(updateUserUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal(err)
				}
				userJson, _ := json.Marshal(response)
				assert.JSONEq(t, tt.expectedBody, string(userJson))

				var userInDB dtos.UserModel
				tx.First(&userInDB, "email = ?", tt.args.userEmail)

				assert.Equal(t, tt.args.input["name"], userInDB.Name, "User name was not updated on database")
				assert.Equal(t, tt.args.input["device_token"], userInDB.DeviceToken, "User token device was not updated on database")
			}
		})
	}
}
