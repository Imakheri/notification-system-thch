package channels

import (
	"errors"
	"strings"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type EmailNotification struct {
	entities.Notification
	Wrapped        entities.NotificationComponent
	RecipientEmail string
}

func (e EmailNotification) Validate() error {
	if len(e.RecipientEmail) <= 0 {
		return errors.New("recipient email cannot be empty")
	}
	if !strings.Contains(e.RecipientEmail, "@") || !strings.Contains(e.RecipientEmail, ".") {
		return errors.New("invalid email structure")
	}
	return nil
}
