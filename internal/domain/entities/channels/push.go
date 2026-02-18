package channels

import (
	"errors"
	"fmt"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type PushStrategy struct{}

func NewPushStrategy() entities.NotificationStrategy {
	return &PushStrategy{}
}

func (ps *PushStrategy) validate(user entities.User) error {
	if len(user.DeviceToken) < 0 {
		return errors.New("user must have a device token")
	}
	return nil
}

func (ps *PushStrategy) Send(user entities.User, notification entities.Notification) error {
	err := ps.validate(user)
	if err != nil {
		return err
	}
	fmt.Println("sending via push")
	return nil
}
