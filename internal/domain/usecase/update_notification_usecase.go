package usecase

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type UpdateNotificationUseCase interface {
	Exec(userID uint, notificationID int, notification entities.Notification) (entities.Notification, error)
}

type updateNotificationUseCase struct {
	updateNotificationRepository gateway.NotificationRepository
}

func NewUpdateNotificationUseCase(updateNotificationRepository gateway.NotificationRepository) UpdateNotificationUseCase {
	return &updateNotificationUseCase{
		updateNotificationRepository: updateNotificationRepository,
	}
}

func (u *updateNotificationUseCase) Exec(userID uint, notificationID int, notification entities.Notification) (entities.Notification, error) {
	notification, err := u.updateNotificationRepository.UpdateNotification(userID, notificationID, notification)
	if err != nil {
		return entities.Notification{}, err
	}
	return notification, nil
}
