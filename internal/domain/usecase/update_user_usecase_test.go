package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
)

func Test_updateUserUseCase_Exec(t *testing.T) {
	type fields struct {
		userRepository func(ctrl *gomock.Controller) gateway.UserRepository
	}
	type args struct {
		userEmail string
		user      entities.User
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
			name: "User updated successfully",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{
						ID:            3,
						Name:          "Fernando Pessoa",
						Password:      "Mensagem123.",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "",
						Notifications: nil,
					}, nil)
					m.EXPECT().UpdateUser(entities.User{
						ID:          3,
						Name:        "Ricardo Reis",
						DeviceToken: "12345qwerty67890",
					}).Return(entities.User{
						ID:          3,
						Name:        "Ricardo Reis",
						Email:       "fernandop@exmaple.com",
						Phone:       "0129834765",
						DeviceToken: "12345qwerty67890",
					}, nil)
					return m
				},
			},
			args: args{
				userEmail: "fernandop@example.com",
				user: entities.User{
					Name:        "Ricardo Reis",
					DeviceToken: "12345qwerty67890",
				},
			},
			want: entities.User{
				ID:            3,
				Name:          "Ricardo Reis",
				Email:         "fernandop@exmaple.com",
				Phone:         "0129834765",
				DeviceToken:   "12345qwerty67890",
				Notifications: nil,
			},
			wantErr: false,
		},
		{
			name: "User do not exists",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{}, errors.New("user not found"))
					return m
				},
			},
			args: args{
				userEmail: "fernandop@example.com",
				user: entities.User{
					Name:        "Ricardo Reis",
					Phone:       "0129834765",
					DeviceToken: "12345qwerty67890",
				},
			},
			want:           entities.User{},
			wantErr:        true,
			ExpectedErrMsg: "user not found",
		},
		{
			name: "Error from database updating user",
			fields: fields{
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{}, context.DeadlineExceeded)
					return m
				},
			},
			args: args{
				userEmail: "fernandop@example.com",
				user: entities.User{
					Name:        "Ricardo Reis",
					Phone:       "0129834765",
					DeviceToken: "12345qwerty67890",
				},
			},
			want:           entities.User{},
			wantErr:        true,
			ExpectedErrMsg: "context deadline exceeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := tt.fields.userRepository(ctrl)

			usecase := NewUpdateUserUseCase(repository)

			got, err := usecase.Exec(tt.args.userEmail, tt.args.user)
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
