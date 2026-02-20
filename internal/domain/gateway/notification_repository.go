package gateway

import (
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type NotificationRepository interface {
	CreateNotification(uint, entities.Notification) (entities.Notification, error)
	GetNotificationsByUser(uint) ([]entities.Notification, error)
	DoesNotificationExistsAndBelongsToUser(uint, uint) (entities.Notification, error)
	UpdateNotification(notification entities.Notification) (entities.Notification, error)
	DeleteNotificationByID(notificationID uint) (uint, error)
	SetSentAt(entities.Notification, time.Time) (entities.Notification, error)
}
