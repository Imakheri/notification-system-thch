package usecase

import (
	"errors"
	"testing"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_getUserUseCase_Exec(t *testing.T) {
	type fields struct {
		userRepository func(ctrl *gomock.Controller) gateway.UserRepository
		jwtService     func(ctrl *gomock.Controller) gateway.JwTokenService
	}
	type args struct {
		userInput entities.User
	}
	type want struct {
		user           entities.User
		token          string
		wantErr        bool
		ExpectedErrMsg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Got user successfully (Login success)",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					passHashed, _ := bcrypt.GenerateFromPassword([]byte("Hamlet123."), 10)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            2,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Password:      string(passHashed),
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					return m
				},
				jwtService: func(ctrl *gomock.Controller) gateway.JwTokenService {
					m := mocks.NewMockJwTokenService(ctrl)
					const expectedToken = "!2#4%6/8)0.qwerty_this_is_a_valid_token-.signature_test=9(7&5$31"
					m.EXPECT().GenerateToken("williams@example.com", uint(2)).Return(expectedToken, nil)
					return m
				},
			},
			args: args{
				userInput: entities.User{
					Password: "Hamlet123.",
					Email:    "williams@example.com",
				},
			},
			want: want{
				user: entities.User{
					ID:            2,
					Name:          "William Shakespeare",
					Email:         "williams@example.com",
					Phone:         "1234567890",
					DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
					Notifications: nil,
				},
				token:   "!2#4%6/8)0.qwerty_this_is_a_valid_token-.signature_test=9(7&5$31",
				wantErr: false,
			},
		},
		{
			name: "Wrong password (Login unsuccess)",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					passHashed, _ := bcrypt.GenerateFromPassword([]byte("Hamlet123."), 10)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            2,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Password:      string(passHashed),
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					return m
				},
				jwtService: func(ctrl *gomock.Controller) gateway.JwTokenService {
					m := mocks.NewMockJwTokenService(ctrl)
					return m
				},
			},
			args: args{
				userInput: entities.User{
					Password: "Isnotmypassword123.",
					Email:    "williams@example.com",
				},
			},
			want: want{
				user:           entities.User{},
				token:          "",
				wantErr:        true,
				ExpectedErrMsg: "the e-mail address or password is incorrect",
			},
		},
		{
			name: "User does not exists",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("example@example.com").Return(entities.User{}, errors.New("user not found"))
					return m
				},
				jwtService: func(ctrl *gomock.Controller) gateway.JwTokenService {
					m := mocks.NewMockJwTokenService(ctrl)
					return m
				},
			},
			args: args{
				userInput: entities.User{
					Email:    "example@example.com",
					Password: "IsAValidPass123!",
				},
			},
			want: want{
				user:           entities.User{},
				token:          "",
				wantErr:        true,
				ExpectedErrMsg: "user not found",
			},
		},
		{
			name: "Got user successfully (Login unsuccess)",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					passHashed, _ := bcrypt.GenerateFromPassword([]byte("Hamlet123."), 10)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            2,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Password:      string(passHashed),
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					return m
				},
				jwtService: func(ctrl *gomock.Controller) gateway.JwTokenService {
					m := mocks.NewMockJwTokenService(ctrl)
					m.EXPECT().GenerateToken("williams@example.com", uint(2)).Return("", errors.New("could not generate JWT"))
					return m
				},
			},
			args: args{
				userInput: entities.User{
					Password: "Hamlet123.",
					Email:    "williams@example.com",
				},
			},
			want: want{
				user:           entities.User{},
				token:          "",
				wantErr:        true,
				ExpectedErrMsg: "could not generate JWT",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepository := tt.fields.userRepository(ctrl)
			jwtService := tt.fields.jwtService(ctrl)

			usecase := NewGetUserUseCase(userRepository, jwtService)

			gotUser, gotToken, err := usecase.Exec(tt.args.userInput)

			if (err != nil) != tt.want.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.want.wantErr)
				return
			}

			if gotUser.ID != tt.want.user.ID ||
				gotUser.Name != tt.want.user.Name ||
				gotUser.Email != tt.want.user.Email ||
				gotUser.Phone != tt.want.user.Phone ||
				gotUser.DeviceToken != tt.want.user.DeviceToken {
				t.Errorf("Exec() gotUser = %v, want %v", gotUser, tt.want.user)
			}

			if gotToken != tt.want.token {
				t.Errorf("Exec() gotToken = %v, want %v", gotToken, tt.want.token)
			}

		})
	}
}
