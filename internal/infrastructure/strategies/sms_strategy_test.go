package strategies

import (
	"net/http"
	"testing"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestSMSStrategy_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockSimulatedApiService(ctrl)
	strategy := NewSMSStrategy(m)

	tests := []struct {
		name           string
		recipient      entities.User
		notification   entities.Notification
		mockBehavior   func()
		expectedStatus int
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "SMS sent successfully on first attempt",
			recipient: entities.User{
				Phone: "1234567890",
			},
			notification: entities.Notification{
				Title:   "Test Notification",
				Content: "This is a test notification",
			},
			mockBehavior: func() {
				m.EXPECT().RandomizeHTTPStatus().Return(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "Invalid phone number structure",
			recipient: entities.User{
				Phone: "1234567",
			},
			notification: entities.Notification{
				Title:   "Test Notification",
				Content: "This is a test notification",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
			ExpectedErrMsg: "invalid phone number",
		},
		{
			name: "Invalid title and content length",
			recipient: entities.User{
				Phone: "1234567890",
			},
			notification: entities.Notification{
				Title:   "Lorem ipsum",
				Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
			ExpectedErrMsg: "in total the content and title should not not exceed 160 characters",
		},
		{
			name: "SMS not sent after 3 attempts",
			recipient: entities.User{
				Phone: "1234567890",
			},
			notification: entities.Notification{
				Title:   "Test Notification",
				Content: "This is a test notification",
			},
			mockBehavior: func() {
				m.EXPECT().RandomizeHTTPStatus().Return(http.StatusInternalServerError).Times(3)
			},
			expectedStatus: http.StatusInternalServerError,
			wantErr:        true,
			ExpectedErrMsg: "an error occurred while trying to send notification via sms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			status, err := strategy.Send("sender", tt.recipient, tt.notification)

			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("Send() error = [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if status != tt.expectedStatus {
				t.Errorf("Send() status = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}
