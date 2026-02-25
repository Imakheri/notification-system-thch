package strategies

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type SMSStrategy struct {
	simulatedApiService gateway.SimulatedApiService
}

func NewSMSStrategy(simulatedApiService gateway.SimulatedApiService) entities.NotificationStrategy {
	return &SMSStrategy{
		simulatedApiService: simulatedApiService,
	}
}

func (ss *SMSStrategy) validate(user entities.User, notification entities.Notification) error {
	if len(user.Phone) <= 0 {
		return errors.New("user must have a phone number")
	}
	if len(user.Phone) < 10 {
		return errors.New("invalid phone number")
	}
	_, err := strconv.Atoi(user.Phone)
	if err != nil {
		return err
	}
	if (len(notification.Content) + len(notification.Title)) > 160 {
		return errors.New("in total the content and title should not not exceed 160 characters")
	}
	return nil
}

func (ss *SMSStrategy) Send(sender string, recipient entities.User, notification entities.Notification) (int, error) {
	err := ss.validate(recipient, notification)
	if err != nil {
		return 0, nil
	}
	maxRetries := 3
	var status int
	for i := 0; i < maxRetries; i++ {
		status = ss.simulatedApiService.RandomizeHTTPStatus()
		if status == http.StatusOK || status == http.StatusCreated {
			return status, nil
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return status, errors.New("an error occurred while trying to send notification via sms")
}
