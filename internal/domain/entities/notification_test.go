package entities

import (
	"reflect"
	"testing"
)

func TestNewNotification(t *testing.T) {
	type args struct {
		title      string
		content    string
		channelID  uint
		recipients []User
	}
	tests := []struct {
		name           string
		args           args
		want           Notification
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Valid New Notification",
			args: args{
				title:     "Test Notification",
				content:   "This is a test notification",
				channelID: 2,
				recipients: []User{
					{
						Email: "example@example.com",
					},
				},
			},
			want: Notification{
				ID:        0,
				CreatedBy: 0,
				SentAt:    nil,
				Title:     "Test Notification",
				Content:   "This is a test notification",
				ChannelID: 2,
				Recipients: []User{
					{
						Email: "example@example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid New Notification (Invalid Recipients)",
			args: args{
				title:      "Test Notification",
				content:    "This is a test notification",
				channelID:  2,
				recipients: nil,
			},
			want:           Notification{},
			wantErr:        true,
			ExpectedErrMsg: "recipients must not be greater than 0",
		},
		{
			name: "Invalid New Notification (Invalid Title)",
			args: args{
				title:      "!1234 #$ 567%&",
				content:    "This is a test notification",
				channelID:  2,
				recipients: nil,
			},
			want:           Notification{},
			wantErr:        true,
			ExpectedErrMsg: "title notification must not have more than 6 digits or symbols",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNotification(tt.args.title, tt.args.content, tt.args.channelID, tt.args.recipients)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewNotification() error = [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotification() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateNotification(t *testing.T) {
	type args struct {
		title      string
		content    string
		channelID  uint
		recipients []User
	}
	tests := []struct {
		name           string
		args           args
		want           Notification
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Invalid New Notification (Invalid Channel)",
			args: args{
				title:      "Test Notification",
				content:    "This is a test notification",
				channelID:  4,
				recipients: nil,
			},
			want:           Notification{},
			wantErr:        true,
			ExpectedErrMsg: "must use a valid channel",
		},
		{
			name: "Invalid New Notification (Invalid Title)",
			args: args{
				title:      "¡Title!",
				content:    "This is a test notification",
				channelID:  4,
				recipients: nil,
			},
			want:           Notification{},
			wantErr:        true,
			ExpectedErrMsg: "title notification must have at least 10 characters",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateNotification(tt.args.title, tt.args.content, tt.args.channelID, tt.args.recipients)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewNotification() error = [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateNotification() got = %v, want %v", got, tt.want)
			}
		})
	}
}
