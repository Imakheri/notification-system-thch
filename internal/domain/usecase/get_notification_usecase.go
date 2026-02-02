package usecase

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type GetNotificationsByUserUseCase interface {
	Exec(userID uint) ([]entities.Notification, error)
}

type getNotificationsByUserUseCase struct {
	notificationRepository gateway.NotificationRepository
}

func NewGetNotificationsByUserUseCase(notificationRepository gateway.NotificationRepository) GetNotificationsByUserUseCase {
	return &getNotificationsByUserUseCase{
		notificationRepository: notificationRepository,
	}
}

func (gnu *getNotificationsByUserUseCase) Exec(userID uint) ([]entities.Notification, error) {
	notifications, err := gnu.notificationRepository.GetNotificationsByUser(userID)
	if err != nil {
		return []entities.Notification{}, err
	}
	return notifications, nil
}
