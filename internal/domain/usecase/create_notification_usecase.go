package usecase

import (
	"errors"
	"fmt"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/entities/channels"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type CreateNotificationUseCase interface {
	Exec(userID uint, userEmail string, notification entities.Notification) (entities.Notification, error)
}

type createNotificationUseCase struct {
	userRepository         gateway.UserRepository
	notificationRepository gateway.NotificationRepository
}

func NewCreateNotificationUseCase(notificationRepository gateway.NotificationRepository, userRepository gateway.UserRepository) CreateNotificationUseCase {
	return &createNotificationUseCase{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
	}
}

func (cnu *createNotificationUseCase) Exec(userID uint, userEmail string, notification entities.Notification) (entities.Notification, error) {
	err := cnu.checkNotificationProperties(notification)
	if err != nil {
		return entities.Notification{}, err
	}

	notification, err = cnu.notificationRepository.CreateNotification(userID, notification)
	if err != nil {
		return entities.Notification{}, err
	}

	for _, recipient := range notification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			_, err := cnu.notificationRepository.DeleteNotificationByID(notification.ID)
			return entities.Notification{}, err
		}
		for _, channel := range notification.Channels {
			currentStrategy := strategySelector(channel.ID)
			err = currentStrategy.Send(user, notification)
			if err != nil {
				_, err := cnu.notificationRepository.DeleteNotificationByID(notification.ID)
				return entities.Notification{}, err
			}
		}
	}

	return notification, nil
}

func (cnu *createNotificationUseCase) checkNotificationProperties(notification entities.Notification) error {
	if len(notification.Title) == 0 {
		return errors.New("notification must have a title")
	}
	if len(notification.Content) == 0 {
		return errors.New("notification must have content")
	}
	if len(notification.Recipients) < 0 {
		return errors.New("recipients must not be greater than 0")
	}

	for i, recipient := range notification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			return errors.New("recipient does not exist")
		}
		notification.Recipients[i].ID = user.ID
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

func strategySelector(idStrategy uint) entities.NotificationStrategy {
	var selectedStrategy entities.NotificationStrategy
	switch idStrategy {
	case 1:
		selectedStrategy = channels.NewEmailStrategy()
	case 2:
		selectedStrategy = channels.NewSMSStrategy()
	case 3:
		selectedStrategy = channels.NewPushStrategy()
	default:
		fmt.Println("Unknown strategy")
	}
	return selectedStrategy
}
