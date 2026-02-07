package usecase

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type CreateNotificationUseCase interface {
	Exec(userID uint, notification entities.Notification) (entities.Notification, error)
}

type createNotificationUseCase struct {
	repository gateway.NotificationRepository
}

func NewCreateNotificationUseCase(repository gateway.NotificationRepository) CreateNotificationUseCase {
	return &createNotificationUseCase{
		repository: repository,
	}
}

func (cnu createNotificationUseCase) Exec(userID uint, notification entities.Notification) (entities.Notification, error) {
	err := checkNotificationProperties(notification)
	if err != nil {
		return entities.Notification{}, err
	}

	notification, err = cnu.repository.CreateNotification(userID, notification)
	if err != nil {
		return entities.Notification{}, err
	}

	return notification, nil
}

func checkNotificationProperties(notification entities.Notification) error {
	if len(notification.Title) == 0 {
		return errors.New("notification must have a title")
	}
	if len(notification.Content) == 0 {
		return errors.New("notification must have content")
	}
	if len(notification.Recipients) < 0 {
		return errors.New("recipients must not be greater than 0")
	}
	if len(notification.Channels) < 0 {
		return errors.New("notification must use at least one channel")
	}
	for _, channel := range notification.Channels {
		if channel.ID == 0 || channel.ID > 3 {
			return errors.New("mush enter a valid channel")
		}
	}
	return nil
}
