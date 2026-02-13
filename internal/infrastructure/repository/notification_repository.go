package repository

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db_connection *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) gateway.NotificationRepository {
	return &NotificationRepository{
		db_connection: db,
	}
}

func (nr *NotificationRepository) CreateNotification(userID uint, notification entities.Notification) (entities.Notification, error) {
	var user entities.User
	result := Database().First(&user, userID)
	if result.Error != nil {
		return entities.Notification{}, errors.New("user not found")
	}

	notification.CreatedBy = userID
	result = Database().Create(&notification)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}

func (nr *NotificationRepository) GetNotificationsByUser(userID uint) ([]entities.Notification, error) {
	var notifications []entities.Notification
	result := Database().Find(&notifications, "created_by = ?", userID)
	if result.Error != nil {
		return []entities.Notification{}, result.Error
	}
	return notifications, nil
}

func (nr *NotificationRepository) DeleteNotificationByID(notificationID uint) (uint, error) {
	result := Database().Delete(&entities.Notification{}, notificationID)
	if result.Error != nil {
		return 0, result.Error
	}
	return notificationID, nil
}

func (nr *NotificationRepository) UpdateNotification(notification entities.Notification) (entities.Notification, error) {
	result := Database().Save(&notification)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}

func (nr *NotificationRepository) DoesNotificationExistsAndBelongsToUser(userID uint, notificationID uint) (entities.Notification, error) {
	var notification entities.Notification
	result := Database().First(&notification, notificationID)
	if result.Error != nil {
		return notification, errors.New("notification not found")
	}
	if notification.CreatedBy != userID {
		return notification, errors.New("notification does not belong to user")
	}
	return notification, nil
}
