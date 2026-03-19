package strategies

import (
	"errors"
	"net/http"
	"testing"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestEmailStrategy_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockSimulatedApiService(ctrl)
	strategy := NewEmailStrategy(m)

	tests := []struct {
		name           string
		recipient      entities.User
		mockBehavior   func()
		expectedStatus int
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name:      "Email sent successfully on first attempt",
			recipient: entities.User{Email: "test@example.com"},
			mockBehavior: func() {
				m.EXPECT().RandomizeHTTPStatus().Return(http.StatusOK, nil)
			},
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Invalid email structure",
			recipient: entities.User{Email: "testexample.com"},
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
			ExpectedErrMsg: "invalid email structure",
		},
		{
			name:      "Email sent  error after 3 attempts",
			recipient: entities.User{Email: "test@example.com"},
			mockBehavior: func() {
				m.EXPECT().RandomizeHTTPStatus().Return(http.StatusInternalServerError, errors.New("an error occurred while trying to send"))
			},
			expectedStatus: http.StatusInternalServerError,
			wantErr:        true,
			ExpectedErrMsg: "an error occurred while trying to send notification via email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			status, err := strategy.Send("sender", tt.recipient, entities.Notification{})

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
