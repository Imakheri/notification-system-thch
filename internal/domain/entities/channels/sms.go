package channels

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type SMSNotification struct {
	entities.Notification
	Wrapped        entities.NotificationComponent
	RecipientPhone string
}

func (s SMSNotification) Validate() error {
	if len(s.RecipientPhone) <= 0 {
		return errors.New("recipient has to have registered a phone number")
	}
	if len(s.RecipientPhone) < 10 {
		return errors.New("recipient phone must have at least 10 digits")
	}
	return nil
}
