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

func Test_updateNotificationUseCase_Exec(t *testing.T) {
	type fields struct {
		updateNotificationRepository func(ctrl *gomock.Controller) gateway.NotificationRepository
		userRepository               func(ctrl *gomock.Controller) gateway.UserRepository
	}
	type args struct {
		userID         uint
		notificationID int
		notification   entities.Notification
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           entities.Notification
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Notification updated successfully",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  1,
						SentAt:     nil,
						Title:      "Test notification #100",
						Content:    "This is a tests notification content",
						ChannelID:  1,
						Recipients: nil,
					}, nil)
					m.EXPECT().UpdateNotification(entities.Notification{
						ID:        100,
						CreatedBy: 0,
						SentAt:    nil,
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly. ",
						ChannelID: 3,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
							{
								ID:    4,
								Email: "fernandop@example.com",
							},
						},
					}).Return(entities.Notification{
						ID:        100,
						CreatedBy: 1,
						SentAt:    nil,
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly. ",
						ChannelID: 3,
						Recipients: []entities.User{
							{
								ID:            3,
								Name:          "Charles Dickens",
								Email:         "charlesd@example.com",
								Phone:         "0987654321",
								DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
								Notifications: nil,
							},
							{
								ID:            4,
								Name:          "Fernando Pessoa",
								Email:         "fernandop@example.com",
								Phone:         "0129834765",
								DeviceToken:   "12345qwerty67890",
								Notifications: nil,
							},
						},
					}, nil)
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("charlesd@example.com").Return(entities.User{
						ID:            3,
						Name:          "Charles Dickens",
						Email:         "charlesd@example.com",
						Phone:         "0987654321",
						DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
						Notifications: nil,
					}, nil)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{
						ID:            4,
						Name:          "Fernando Pessoa",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					}, nil)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
				notification: entities.Notification{
					Title:     "Earthquake detected in your area!",
					Content:   "Drop, cover, and hold on. Please proceed calmly. ",
					ChannelID: 3,
					Recipients: []entities.User{
						{
							Email: "charlesd@example.com",
						},
						{
							Email: "fernandop@example.com",
						},
					},
				},
			},
			want: entities.Notification{
				ID:        100,
				CreatedBy: 1,
				SentAt:    nil,
				Title:     "Earthquake detected in your area!",
				Content:   "Drop, cover, and hold on. Please proceed calmly. ",
				ChannelID: 3,
				Recipients: []entities.User{
					{
						ID:            3,
						Name:          "Charles Dickens",
						Email:         "charlesd@example.com",
						Phone:         "0987654321",
						DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
						Notifications: nil,
					},
					{
						ID:            4,
						Name:          "Fernando Pessoa",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Notification does not exist",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{}, errors.New("notification not found"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
				notification:   entities.Notification{},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "notification not found",
		},
		{
			name: "Notification does not belong to user",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{}, errors.New("notification does not belong to user"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
				notification:   entities.Notification{},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "notification does not belong to user",
		},
		{
			name: "Invalid recipient on notification",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(2), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  2,
						SentAt:     nil,
						Title:      "Test notification #100",
						Content:    "This is a tests notification content",
						ChannelID:  1,
						Recipients: nil,
					}, nil)
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            2,
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
				userID:         2,
				notificationID: 100,
				notification: entities.Notification{
					Title:     "Earthquake detected in your area!",
					Content:   "Drop, cover, and hold on. Please proceed calmly. ",
					ChannelID: 3,
					Recipients: []entities.User{
						{
							Email: "williams@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "invalid recipient",
		},
		{
			name: "Non existent recipient on notification",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  1,
						SentAt:     nil,
						Title:      "Test notification #100",
						Content:    "This is a tests notification content",
						ChannelID:  1,
						Recipients: nil,
					}, nil)
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{
						ID:            2,
						Name:          "William Shakespeare",
						Email:         "williams@example.com",
						Phone:         "1234567890",
						DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
						Notifications: nil,
					}, nil)
					m.EXPECT().GetUserByEmail("example@example.com").Return(entities.User{}, errors.New("recipient does not exist"))
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
				notification: entities.Notification{
					Title:     "Earthquake detected in your area!",
					Content:   "Drop, cover, and hold on. Please proceed calmly. ",
					ChannelID: 3,
					Recipients: []entities.User{
						{
							Email: "williams@example.com",
						},
						{
							Email: "example@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "recipient does not exist",
		},
		{
			name: "Notification updated successfully",
			fields: fields{
				updateNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  1,
						SentAt:     nil,
						Title:      "Test notification #100",
						Content:    "This is a tests notification content",
						ChannelID:  1,
						Recipients: nil,
					}, nil)
					m.EXPECT().UpdateNotification(entities.Notification{
						ID:        100,
						CreatedBy: 0,
						SentAt:    nil,
						Title:     "Earthquake detected in your area!",
						Content:   "Drop, cover, and hold on. Please proceed calmly. ",
						ChannelID: 3,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
							{
								ID:    4,
								Email: "fernandop@example.com",
							},
						},
					}).Return(entities.Notification{}, errors.New("db error: connection lost during update"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("charlesd@example.com").Return(entities.User{
						ID:            3,
						Name:          "Charles Dickens",
						Email:         "charlesd@example.com",
						Phone:         "0987654321",
						DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
						Notifications: nil,
					}, nil)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{
						ID:            4,
						Name:          "Fernando Pessoa",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					}, nil)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
				notification: entities.Notification{
					Title:     "Earthquake detected in your area!",
					Content:   "Drop, cover, and hold on. Please proceed calmly. ",
					ChannelID: 3,
					Recipients: []entities.User{
						{
							Email: "charlesd@example.com",
						},
						{
							Email: "fernandop@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "db error: connection lost during update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			notificationRepository := tt.fields.updateNotificationRepository(ctrl)
			userRepository := tt.fields.userRepository(ctrl)
			usecase := NewUpdateNotificationUseCase(notificationRepository, userRepository)
			got, err := usecase.Exec(tt.args.userID, tt.args.notificationID, tt.args.notification)
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
