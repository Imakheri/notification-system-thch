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
	notificationRepository gateway.NotificationRepository
	userRepository         gateway.UserRepository
}

func NewUpdateNotificationUseCase(updateNotificationRepository gateway.NotificationRepository, userRepository gateway.UserRepository) UpdateNotificationUseCase {
	return &updateNotificationUseCase{
		notificationRepository: updateNotificationRepository,
		userRepository:         userRepository,
	}
}

func (u *updateNotificationUseCase) Exec(userID uint, notificationID int, notification entities.Notification) (entities.Notification, error) {
	notificationFromDB, err := u.notificationRepository.DoesNotificationExistsAndBelongsToUser(userID, uint(notificationID))
	if err != nil {
		return entities.Notification{}, err
	}

	notification.ID = notificationFromDB.ID

	if len(notification.Recipients) > 0 {
		for i, recipient := range notification.Recipients {
			user, err := u.userRepository.GetUserByEmail(recipient.Email)
			if err != nil {
				return entities.Notification{}, errors.New("recipient does not exist")
			}
			if userID == user.ID {
				return entities.Notification{}, errors.New("invalid recipient")
			}
			notification.Recipients[i].ID = user.ID
		}
	}

	_, err = u.notificationRepository.UpdateNotification(notification)
	if err != nil {
		return entities.Notification{}, err
	}

	updatedNotificationFromDB, err := u.notificationRepository.DoesNotificationExistsAndBelongsToUser(userID, uint(notificationID))
	if err != nil {
		return entities.Notification{}, err
	}

	return updatedNotificationFromDB, nil
}
