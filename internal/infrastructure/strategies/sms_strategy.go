package strategies

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	if len(user.Phone) == 0 {
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
		return http.StatusBadRequest, err
	}
	status, err := ss.simulatedApiService.RandomizeHTTPStatus()
	if err != nil {
		return status, errors.New(fmt.Sprintf("%v %v", err, "notification via sms"))
	}
	return status, nil
}
