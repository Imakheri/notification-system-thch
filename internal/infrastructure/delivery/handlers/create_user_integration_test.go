//go:build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateUser(t *testing.T) {
	type args struct {
		input gin.H
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "User created successfully",
			args: args{
				input: gin.H{
					"name":         "William Shakespeare",
					"email":        "williams@example.com",
					"password":     "Hamlet123.",
					"phone":        "1234567890",
					"device_token": "0q1e2r3t4y5a6b7c8d9e",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
								"user": {
									"name":         "William Shakespeare",
									"email":        "williams@example.com",
									"phone":        "1234567890",
									"device_token": "0q1e2r3t4y5a6b7c8d9e"
								}
							}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TestDbStruct.DatabaseConnection.Begin()
			defer tx.Rollback()

			dbWrapper := &repository.Database{
				DatabaseConnection: tx,
			}
			createUserRepository := repository.NewUserRepository(dbWrapper)
			createUserUseCase := usecase.NewCreateUserUseCase(createUserRepository)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			jsonInput, _ := json.Marshal(tt.args.input)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonInput))
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := CreateUserHandler(createUserUseCase)
			h(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actualResponse map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &actualResponse)

			var tokenString string
			if userMap, ok := actualResponse["user"].(map[string]interface{}); ok {
				tokenString, _ = userMap["token"].(string)
				delete(userMap, "token")
			}

			actualJSON, _ := json.Marshal(actualResponse)
			assert.JSONEq(t, tt.expectedBody, string(actualJSON))

			var user dtos.UserModel
			result := tx.Where("email = ?", tt.args.input["email"]).First(&user)
			assert.NoError(t, result.Error, "user should exists on database")
			assert.NotEqual(t, tt.args.input["password"], user.Password, "Password was saved as plain text")

			if tokenString != "" {
				token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
				assert.NoError(t, err, "Not a valid JWToken")

				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					assert.Equal(t, tt.args.input["email"], claims["email"])
					assert.Equal(t, float64(user.ID), claims["id"])
				}
			}
		})
	}
}
