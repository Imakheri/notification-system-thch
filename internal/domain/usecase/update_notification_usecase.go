package usecase

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type UpdateNotificationUseCase interface {
	Exec(userID uint, notificationID int, notification entities.Notification) (entities.Notification, error)
}

type updateNotificationUseCase struct {
	updateNotificationRepository gateway.NotificationRepository
	userRepository               gateway.UserRepository
}

func NewUpdateNotificationUseCase(updateNotificationRepository gateway.NotificationRepository, userRepository gateway.UserRepository) UpdateNotificationUseCase {
	return &updateNotificationUseCase{
		updateNotificationRepository: updateNotificationRepository,
		userRepository:               userRepository,
	}
}

func (u *updateNotificationUseCase) Exec(userID uint, notificationID int, notificationDTO entities.Notification) (entities.Notification, error) {
	notification, err := u.updateNotificationRepository.DoesNotificationExistsAndBelongsToUser(userID, uint(notificationID))
	if err != nil {
		return entities.Notification{}, err
	}

	if len(notificationDTO.Title) <= 0 {
		notificationDTO.Title = notification.Title
	} else {
		notification.Title = notificationDTO.Title
	}

	if len(notificationDTO.Content) <= 0 {
		notificationDTO.Content = notification.Content
	} else {
		notification.Content = notificationDTO.Content
	}

	if len(notificationDTO.Channels) <= 0 {
		notificationDTO.Channels = notification.Channels
	} else {
		notification.Channels = notificationDTO.Channels
	}

	if len(notificationDTO.Recipients) <= 0 {
		notificationDTO.Recipients = notification.Recipients
	} else {
		for i, recipient := range notificationDTO.Recipients {
			user, err := u.userRepository.GetUserByEmail(recipient.Email)
			if err != nil {
				return entities.Notification{}, errors.New("recipient does not exist")
			}
			if userID == user.ID {
				return entities.Notification{}, errors.New("invalid recipient")
			}
			notificationDTO.Recipients[i].ID = user.ID
		}
		notification.Recipients = notificationDTO.Recipients
	}

	notification, err = u.updateNotificationRepository.UpdateNotification(notification)
	if err != nil {
		return entities.Notification{}, err
	}
	return notification, nil
}
