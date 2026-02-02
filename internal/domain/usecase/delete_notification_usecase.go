package usecase

import (
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type DeleteNotificationUseCase interface {
	Exec(notificationID int) (int, error)
}

type deleteNotificationUseCase struct {
	deleteNotificationRepository gateway.NotificationRepository
}

func NewDeleteNotification(deleteNotificationRepository gateway.NotificationRepository) DeleteNotificationUseCase {
	return &deleteNotificationUseCase{
		deleteNotificationRepository: deleteNotificationRepository,
	}
}

func (dnu deleteNotificationUseCase) Exec(notificationID int) (int, error) {
	notificationID, err := dnu.deleteNotificationRepository.DeleteNotificationByID(notificationID)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}
