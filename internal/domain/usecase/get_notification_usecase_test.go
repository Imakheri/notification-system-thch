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

func Test_getNotificationsByUserUseCase_Exec(t *testing.T) {
	type fields struct {
		getNotificationRepository func(ctrl *gomock.Controller) gateway.NotificationRepository
	}
	type args struct {
		userID uint
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           []entities.Notification
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Get all notifications successfully",
			fields: fields{
				getNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().GetNotificationsByUser(uint(1)).Return([]entities.Notification{
						{
							ID:        1,
							CreatedBy: 1,
							SentAt:    nil,
							Title:     "Test Notification",
							Content:   "This is a test notification",
							ChannelID: 2,
							Recipients: []entities.User{
								{
									ID:            2,
									Name:          "William Shakespeare",
									Email:         "williams@example.com",
									Phone:         "1234567890",
									DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
									Notifications: nil,
								},
							},
						},
						{
							ID:        2,
							CreatedBy: 1,
							SentAt:    nil,
							Title:     "Test Notification #2",
							Content:   "This is a test notification number 2",
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
						},
					}, nil)
					return m
				},
			},
			args: args{
				userID: 1,
			},
			want: []entities.Notification{
				{
					ID:        1,
					CreatedBy: 1,
					SentAt:    nil,
					Title:     "Test Notification",
					Content:   "This is a test notification",
					ChannelID: 2,
					Recipients: []entities.User{
						{
							ID:            2,
							Name:          "William Shakespeare",
							Email:         "williams@example.com",
							Phone:         "1234567890",
							DeviceToken:   "0q1e2r3t4y5a6b7c8d9e",
							Notifications: nil,
						},
					},
				},
				{
					ID:        2,
					CreatedBy: 1,
					SentAt:    nil,
					Title:     "Test Notification #2",
					Content:   "This is a test notification number 2",
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
				},
			},
			wantErr: false,
		},
		{
			name: "User has no notifications",
			fields: fields{
				getNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().GetNotificationsByUser(uint(1)).Return([]entities.Notification{}, errors.New("user has no notifications"))
					return m
				},
			},
			args: args{
				userID: 1,
			},
			want:           []entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "user has no notifications",
		},
		{
			name: "Database returns data corruption getting notifications",
			fields: fields{
				getNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().
						GetNotificationsByUser(uint(1)).
						Return(nil, errors.New("sql: Scan error on column 'content': incompatible type"))

					return m
				},
			},
			args: args{
				userID: 1,
			},
			want:           []entities.Notification{},
			wantErr:        true,
			ExpectedErrMsg: "sql: Scan error on column 'content': incompatible type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := tt.fields.getNotificationRepository(ctrl)
			usecase := NewGetNotificationsByUserUseCase(repository)
			got, err := usecase.Exec(tt.args.userID)
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
