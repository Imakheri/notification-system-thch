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

func Test_deleteNotificationUseCase_Exec(t *testing.T) {
	type fields struct {
		deleteNotificationRepository func(ctrl *gomock.Controller) gateway.NotificationRepository
	}
	type args struct {
		userID         uint
		notificationID uint
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           uint
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Notification deleted successfully",
			fields: fields{
				deleteNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  1,
						SentAt:     nil,
						Title:      "Test Notification",
						Content:    "This is a test notification",
						ChannelID:  2,
						Recipients: nil,
					}, nil)
					m.EXPECT().DeleteNotificationByID(uint(100)).Return(uint(100), nil)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "Error from database deleting notification",
			fields: fields{
				deleteNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(100)).Return(entities.Notification{
						ID:         100,
						CreatedBy:  1,
						SentAt:     nil,
						Title:      "Test Notification",
						Content:    "This is a test notification",
						ChannelID:  2,
						Recipients: nil,
					}, nil)
					m.EXPECT().DeleteNotificationByID(uint(100)).Return(uint(0), context.DeadlineExceeded)
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 100,
			},
			want:           0,
			wantErr:        true,
			ExpectedErrMsg: "context deadline exceeded",
		},
		{
			name: "Notification does not exists (Invalid notification ID)",
			fields: fields{
				deleteNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(0)).Return(entities.Notification{}, errors.New("notification not found"))
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 0,
			},
			want:           0,
			wantErr:        true,
			ExpectedErrMsg: "notification not found",
		},
		{
			name: "Notification exists but not belongs to user",
			fields: fields{
				deleteNotificationRepository: func(ctrl *gomock.Controller) gateway.NotificationRepository {
					m := mocks.NewMockNotificationRepository(ctrl)
					m.EXPECT().DoesNotificationExistsAndBelongsToUser(uint(1), uint(25)).Return(entities.Notification{
						ID:         25,
						CreatedBy:  2,
						SentAt:     nil,
						Title:      "Test Notification",
						Content:    "This is a test notification",
						ChannelID:  1,
						Recipients: nil,
					}, errors.New("notification does not belong to user"))
					return m
				},
			},
			args: args{
				userID:         1,
				notificationID: 25,
			},
			want:           0,
			wantErr:        true,
			ExpectedErrMsg: "notification does not belong to user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := tt.fields.deleteNotificationRepository(ctrl)

			usecase := NewDeleteNotificationUseCase(repository)

			got, err := usecase.Exec(tt.args.userID, tt.args.notificationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewNotification() error = [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if got != tt.want {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
