package channels

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type SMSStrategy struct{}

func NewSMSStrategy() entities.NotificationStrategy {
	return &SMSStrategy{}
}

func (ss *SMSStrategy) validate(user entities.User) error {
	if len(user.Phone) <= 0 {
		return errors.New("user must have a phone number")
	}
	if len(user.Phone) < 11 {
		return errors.New("invalid phone number")
	}
	_, err := strconv.Atoi(user.Phone)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SMSStrategy) Send(user entities.User, notification entities.Notification) error {
	err := ss.validate(user)
	if err != nil {
		return err
	}
	fmt.Println("sending via sms")
	return nil
}
