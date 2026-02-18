package channels

import (
	"errors"
	"fmt"
	"strings"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type EmailStrategy struct{}

func NewEmailStrategy() entities.NotificationStrategy {
	return &EmailStrategy{}
}

func (es *EmailStrategy) validate(user entities.User) error {
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		return errors.New("invalid email structure")
	}
	return nil
}

func (es *EmailStrategy) Send(user entities.User, notification entities.Notification) error {
	err := es.validate(user)
	if err != nil {
		return err
	}
	fmt.Println("sending via email")
	return nil
}
