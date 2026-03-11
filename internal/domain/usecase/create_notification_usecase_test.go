package usecase

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
)

func Test_createNotificationUseCase_Exec(t *testing.T) {
	var fixedTime = time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	type fields struct {
		userRepository         func(ctrl *gomock.Controller) gateway.UserRepository
		notificationRepository func(ctrl *gomock.Controller) gateway.NotificationRepository
		channelRepository      func(ctrl *gomock.Controller) gateway.ChannelRepository
		simulatedApiService    func(ctrl *gomock.Controller) gateway.SimulatedApiService
		clock                  func(ctrl *gomock.Controller) gateway.Clock
	}
	type args struct {
		userID       uint
		userEmail    string
		notification entities.Notification
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
			name: "Notification created and sent successfully",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(3)).Return(entities.Channel{ID: 3, Name: "push"}, nil)
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
					}, nil).Times(2)
					m.EXPECT().GetUserByEmail("fernandop@example.com").Return(entities.User{
						ID:            4,
						Name:          "Fernando Pessoa",
						Password:      "Mensagem123.",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					}, nil).Times(2)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().CreateNotification(uint(2), entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
								Password:      "Mensagem123.",
								Email:         "fernandop@example.com",
								Phone:         "0129834765",
								DeviceToken:   "12345qwerty67890",
								Notifications: nil,
							},
						},
					}, nil)
					m.EXPECT().SetSentAt(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
								Password:      "Mensagem123.",
								Email:         "fernandop@example.com",
								Phone:         "0129834765",
								DeviceToken:   "12345qwerty67890",
								Notifications: nil,
							},
						},
					}, fixedTime).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    &fixedTime,
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
								Password:      "Mensagem123.",
								Email:         "fernandop@example.com",
								Phone:         "0129834765",
								DeviceToken:   "12345qwerty67890",
								Notifications: nil,
							},
						},
					}, nil)
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					m.EXPECT().RandomizeHTTPStatus().Return(http.StatusInternalServerError).Times(2)
					m.EXPECT().RandomizeHTTPStatus().Return(http.StatusOK).Times(2)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					m.EXPECT().Now().Return(fixedTime).AnyTimes()
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
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
				ID:        1,
				CreatedBy: 2,
				SentAt:    &fixedTime,
				Title:     "Test Notification",
				Content:   "This is a test notification",
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
						Password:      "Mensagem123.",
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
			name: "Error creating notification",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(3)).Return(entities.Channel{ID: 3, Name: "push"}, nil)
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
						Password:      "Mensagem123.",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					}, nil)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().CreateNotification(uint(2), entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
					}).Return(entities.Notification{}, errors.New("driver: bad connection"))
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
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
			ExpectedErrMsg: "driver: bad connection",
		},
		{
			name: "Error sending notification",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(2)).Return(entities.Channel{ID: 2, Name: "sms"}, nil)
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
					}, nil).Times(2)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().CreateNotification(uint(2), entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 2,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
						},
					}).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 2,
						Recipients: []entities.User{
							{
								ID:            3,
								Name:          "Charles Dickens",
								Email:         "charlesd@example.com",
								Phone:         "0987654321",
								DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
								Notifications: nil,
							},
						},
					}, nil)
					m.EXPECT().DeleteNotificationByID(uint(1)).Return(uint(1), nil)
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					m.EXPECT().RandomizeHTTPStatus().Return(http.StatusInternalServerError).Times(3)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 2,
					Recipients: []entities.User{
						{
							Email: "charlesd@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "an error occurred while trying to send notification via sms",
		},
		{
			name: "Error deleting created notification",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(3)).Return(entities.Channel{ID: 3, Name: "push"}, nil)
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
					}, nil).Times(2)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().CreateNotification(uint(2), entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 3,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
						},
					}).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
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
						},
					}, nil)
					m.EXPECT().DeleteNotificationByID(uint(1)).Return(uint(0), errors.New("sql: no rows in result set"))
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					m.EXPECT().RandomizeHTTPStatus().Return(http.StatusInternalServerError).Times(3)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 3,
					Recipients: []entities.User{
						{
							Email: "charlesd@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "sql: no rows in result set",
		},
		{
			name: "Error updating sent_at when notification was sent",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(1)).Return(entities.Channel{ID: 1, Name: "email"}, nil)
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
					}, nil).Times(2)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().CreateNotification(uint(2), entities.Notification{
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								ID:    3,
								Email: "charlesd@example.com",
							},
						},
					}).Return(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								ID:            3,
								Name:          "Charles Dickens",
								Email:         "charlesd@example.com",
								Phone:         "0987654321",
								DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
								Notifications: nil,
							},
						},
					}, nil)
					m.EXPECT().SetSentAt(entities.Notification{
						ID:        1,
						CreatedBy: 2,
						SentAt:    nil,
						Title:     "Test Notification",
						Content:   "This is a test notification",
						ChannelID: 1,
						Recipients: []entities.User{
							{
								ID:            3,
								Name:          "Charles Dickens",
								Email:         "charlesd@example.com",
								Phone:         "0987654321",
								DeviceToken:   "e0d8c7b6a5y4t3r2e1q0",
								Notifications: nil,
							},
						},
					}, fixedTime).Return(entities.Notification{}, errors.New("db error: connection lost during update"))
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					m.EXPECT().RandomizeHTTPStatus().Return(http.StatusCreated)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					m.EXPECT().Now().Return(fixedTime).AnyTimes()
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 1,
					Recipients: []entities.User{
						{
							Email: "charlesd@example.com",
						},
					},
				},
			},
			want:           entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "db error: connection lost during update",
		},
		{
			name: "Invalid recipient on notification",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(2)).Return(entities.Channel{ID: 2, Name: "sms"}, nil)
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
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 2,
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
			name: "Recipient does not exists on notification",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(2)).Return(entities.Channel{ID: 2, Name: "sms"}, nil)
					return m
				},
				userRepository: func(ctrl *gomock.Controller) gateway.UserRepository {
					m := mocks.NewMockUserRepository(ctrl)
					m.EXPECT().GetUserByEmail("example@example.com").Return(entities.User{}, errors.New("recipient does not exist"))
					return m
				},
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    2,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 2,
					Recipients: []entities.User{
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
			name: "Error getting sender information",
			fields: fields{
				channelRepository: func(ctrl *gomock.Controller) gateway.ChannelRepository {
					m := mocks.NewMockChannelRepository(ctrl)
					m.EXPECT().GetChannel(uint(1)).Return(entities.Channel{ID: 1, Name: "email"}, nil)
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
						Password:      "Mensagem123.",
						Email:         "fernandop@example.com",
						Phone:         "0129834765",
						DeviceToken:   "12345qwerty67890",
						Notifications: nil,
					}, nil)
					m.EXPECT().GetUserByEmail("williams@example.com").Return(entities.User{}, errors.New("cannot get sender information"))
					return m
				},
				notificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					return m
				},
				simulatedApiService: func(ctrl *gomock.Controller) gateway.SimulatedApiService {
					m := mocks.NewMockSimulatedApiService(ctrl)
					return m
				},
				clock: func(ctrl *gomock.Controller) gateway.Clock {
					m := mocks.NewMockClock(ctrl)
					return m
				},
			},
			args: args{
				userID:    0,
				userEmail: "williams@example.com",
				notification: entities.Notification{
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 1,
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
			ExpectedErrMsg: "cannot get sender information",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepository := tt.fields.userRepository(ctrl)
			notificationRepository := tt.fields.notificationRepository(ctrl)
			channelRepository := tt.fields.channelRepository(ctrl)
			simulatedApiService := tt.fields.simulatedApiService(ctrl)
			clock := tt.fields.clock(ctrl)
			usecase := NewCreateNotificationUseCase(notificationRepository, userRepository, channelRepository, simulatedApiService, clock)

			got, err := usecase.Exec(tt.args.userID, tt.args.userEmail, tt.args.notification)
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
