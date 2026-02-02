package gateway

import "github.com/imakheri/notifications-thch/internal/domain/entities"

type NotificationRepository interface {
	CreateNotification(uint, entities.Notification) (entities.Notification, error)
	GetNotificationsByUser(uint) ([]entities.Notification, error)
	UpdateNotification(userID uint, notificationID int, notification entities.Notification) (entities.Notification, error)
	DeleteNotificationByID(notificationID int) (int, error)
}
