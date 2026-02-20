package repository

import (
	"errors"
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type NotificationRepository struct {
	db *Database
}

func NewNotificationRepository(db *Database) gateway.NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (nr *NotificationRepository) CreateNotification(userID uint, notification entities.Notification) (entities.Notification, error) {
	var user entities.User
	result := nr.db.DatabaseConnection.First(&user, userID)
	if result.Error != nil {
		return entities.Notification{}, errors.New("user not found")
	}
	notification.CreatedBy = userID
	result = nr.db.DatabaseConnection.Create(&notification)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}

func (nr *NotificationRepository) GetNotificationsByUser(userID uint) ([]entities.Notification, error) {
	var notifications []entities.Notification
	result := nr.db.DatabaseConnection.Find(&notifications, "created_by = ?", userID)
	if result.Error != nil {
		return []entities.Notification{}, result.Error
	}
	return notifications, nil
}

func (nr *NotificationRepository) DeleteNotificationByID(notificationID uint) (uint, error) {
	result := nr.db.DatabaseConnection.Delete(&entities.Notification{}, notificationID)
	if result.Error != nil {
		return 0, result.Error
	}
	return notificationID, nil
}

func (nr *NotificationRepository) UpdateNotification(notification entities.Notification) (entities.Notification, error) {
	result := nr.db.DatabaseConnection.Save(&notification)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}

func (nr *NotificationRepository) DoesNotificationExistsAndBelongsToUser(userID uint, notificationID uint) (entities.Notification, error) {
	var notification entities.Notification
	result := nr.db.DatabaseConnection.First(&notification, notificationID)
	if result.Error != nil {
		return notification, errors.New("notification not found")
	}
	if notification.CreatedBy != userID {
		return notification, errors.New("notification does not belong to user")
	}
	return notification, nil
}

func (nr *NotificationRepository) SetSentAt(notification entities.Notification, time time.Time) (entities.Notification, error) {
	result := nr.db.DatabaseConnection.Model(&notification).Where("id = ?", notification.ID).Update("sent_at", time)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}
