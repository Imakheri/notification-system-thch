package strategies

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type PushNotificationPayload struct {
	DeviceToken string `json:"device_token"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type PushStrategy struct {
	simulatedApiService gateway.SimulatedApiService
}

func NewPushStrategy(simulatedApiService gateway.SimulatedApiService) entities.NotificationStrategy {
	return &PushStrategy{
		simulatedApiService: simulatedApiService,
	}
}

func (ps *PushStrategy) validate(user entities.User, notification entities.Notification) error {
	if len(user.DeviceToken) == 0 {
		return errors.New("user must have a device token")
	}
	_, err := payloadFormater(user, notification)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PushStrategy) Send(sender string, recipient entities.User, notification entities.Notification) (int, error) {
	err := ps.validate(recipient, notification)
	if err != nil {
		return http.StatusBadRequest, err
	}
	status, err := ps.simulatedApiService.RandomizeHTTPStatus()
	if err != nil {
		return status, errors.New(fmt.Sprintf("%v %v", err, "push notification"))
	}
	return status, nil
}

func payloadFormater(user entities.User, notification entities.Notification) ([]byte, error) {
	pushNotificationPayload := PushNotificationPayload{
		DeviceToken: user.DeviceToken,
		Title:       notification.Title,
		Content:     notification.Content,
	}
	payload, err := json.Marshal(pushNotificationPayload)
	if err != nil {
		return []byte{}, errors.New("an error occurred while trying to serialize push notification")
	}
	return payload, nil
}
