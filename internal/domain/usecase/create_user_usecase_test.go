package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
)

func Test_createUserUseCase_Exec(t *testing.T) {
	type fields struct {
		userRepository func(ctrl *gomock.Controller) gateway.UserRepository
	}
	type args struct {
		user entities.User
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           entities.User
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "User created successfully",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{}, errors.New("user not found"))
					m.EXPECT().CreateUser(gomock.Any()).Return(entities.User{
						ID:            1,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					return m
				},
			},
			args: args{
				user: entities.User{
					Name:        "William Shakespeare",
					Password:    "Hamlet!123.",
					Email:       "williams@example.com",
					Phone:       "1234567890",
					DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
				},
			},
			want: entities.User{
				ID:            1,
				Name:          "William Shakespeare",
				Email:         "williams@example.com",
				Phone:         "1234567890",
				DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
				Notifications: nil,
			},
			wantErr: false,
		},
		{
			name: "User already exists",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            1,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					return m
				},
			},
			args: args{
				user: entities.User{
					Name:        "William Shakespeare",
					Password:    "Hamlet!123.",
					Email:       "williams@example.com",
					Phone:       "1234567890",
					DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
				},
			},
			want:           entities.User{},
			wantErr:        true,
			ExpectedErrMsg: "user already exists",
		},
		{
			name: "Database error creating an user",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{}, errors.New("user not found"))
					m.EXPECT().CreateUser(gomock.Any()).Return(entities.User{}, errors.New("db error: insert failed due to connection loss"))
					return m
				},
			},
			args: args{
				user: entities.User{
					Name:        "William Shakespeare",
					Password:    "Hamlet!123.",
					Email:       "williams@example.com",
					Phone:       "1234567890",
					DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
				},
			},
			want:           entities.User{},
			wantErr:        true,
			ExpectedErrMsg: "db error: insert failed due to connection loss",
		},
		{
			name: "Error on Bcrypt, password is too long",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{}, errors.New("user not found"))
					return m
				},
			},
			args: args{
				user: entities.User{
					Name:        "William Shakespeare",
					Password:    "This_is_a_extremely_long_password_it_exceeds_the_seventy_two_bytes_that_supports_bcrypt_to_force_the_error_12345_.!#",
					Email:       "williams@example.com",
					Phone:       "1234567890",
					DeviceToken: "0q1e2r3t4y5a6b7c8d9e",
				},
			},
			want:           entities.User{},
			wantErr:        true,
			ExpectedErrMsg: "failed to encrypt password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := tt.fields.userRepository(ctrl)

			usecase := NewCreateUserUseCase(repository)

			got, err := usecase.Exec(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewNotification() error = [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
