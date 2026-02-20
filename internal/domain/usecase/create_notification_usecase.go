package usecase

import (
	"errors"
	"log"
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/entities/channels"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type CreateNotificationUseCase interface {
	Exec(userID uint, notification entities.Notification) (entities.Notification, error)
}

type createNotificationUseCase struct {
	userRepository         gateway.UserRepository
	notificationRepository gateway.NotificationRepository
	channelRepository      gateway.ChannelRepository
}

func NewCreateNotificationUseCase(notificationRepository gateway.NotificationRepository, userRepository gateway.UserRepository, channelRepository gateway.ChannelRepository) CreateNotificationUseCase {
	return &createNotificationUseCase{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		channelRepository:      channelRepository,
	}
}

func (cnu *createNotificationUseCase) Exec(userID uint, notification entities.Notification) (entities.Notification, error) {
	notification, err := entities.CheckNotificationProperties(notification)
	if err != nil {
		return entities.Notification{}, err
	}

	channel, err := cnu.channelRepository.GetChannel(notification.ChannelID)
	notification.Channel = channel

	for i, recipient := range notification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			return entities.Notification{}, errors.New("recipient does not exist")
		}
		if userID == user.ID {
			return entities.Notification{}, errors.New("invalid recipient")
		}
		notification.Recipients[i].ID = user.ID
	}

	newNotification, err := cnu.notificationRepository.CreateNotification(userID, notification)
	if err != nil {
		return entities.Notification{}, err
	}

	for _, recipient := range newNotification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			_, err := cnu.notificationRepository.DeleteNotificationByID(newNotification.ID)
			return entities.Notification{}, err
		}
		currentStrategy := strategySelector(notification.Channel.ID)
		err = currentStrategy.Send(user, newNotification)
		if err != nil {
			_, errorDeleting := cnu.notificationRepository.DeleteNotificationByID(newNotification.ID)
			if errorDeleting != nil {
				return entities.Notification{}, errorDeleting
			}
			return entities.Notification{}, err
		}
		newNotification, err = cnu.notificationRepository.SetSentAt(newNotification, time.Now())
		if err != nil {
			return entities.Notification{}, err
		}
	}

	return newNotification, nil
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
		log.Fatal("Unknown strategy")
	}
	return selectedStrategy
}
