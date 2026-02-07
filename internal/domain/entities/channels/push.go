package channels

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type PushNotification struct {
	entities.Notification
	Wrapped     entities.NotificationComponent
	DeviceToken string
	Platform    string
}

func (p *PushNotification) Validate() error {
	if len(p.DeviceToken) <= 0 {
		return errors.New("recipient user has to have a registered device")
	}
	return nil
}
