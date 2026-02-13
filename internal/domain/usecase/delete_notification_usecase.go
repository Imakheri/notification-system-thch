package usecase

import (
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type DeleteNotificationUseCase interface {
	Exec(userID uint, notificationID uint) (uint, error)
}

type deleteNotificationUseCase struct {
	deleteNotificationRepository gateway.NotificationRepository
}

func NewDeleteNotification(deleteNotificationRepository gateway.NotificationRepository) DeleteNotificationUseCase {
	return &deleteNotificationUseCase{
		deleteNotificationRepository: deleteNotificationRepository,
	}
}

func (dnu *deleteNotificationUseCase) Exec(userID uint, notificationID uint) (uint, error) {
	_, err := dnu.deleteNotificationRepository.DoesNotificationExistsAndBelongsToUser(userID, uint(notificationID))
	if err != nil {
		return 0, err
	}
	notificationID, err = dnu.deleteNotificationRepository.DeleteNotificationByID(uint(notificationID))
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}
