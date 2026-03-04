package repository

import (
	"errors"
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
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
	var userModel dtos.UserModel
	result := nr.db.DatabaseConnection.First(&userModel, userID)
	if result.Error != nil {
		return entities.Notification{}, errors.New("user not found")
	}
	notificationModel := dtos.NotificationToModel(notification)
	notificationModel.CreatedBy = userID
	result = nr.db.DatabaseConnection.Create(&notificationModel)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	modelEntity := dtos.NotificationModelToEntity(notificationModel)
	return modelEntity, nil
}

func (nr *NotificationRepository) GetNotificationsByUser(userID uint) ([]entities.Notification, error) {
	var notifications []dtos.NotificationModel
	result := nr.db.DatabaseConnection.Preload("Recipients").Find(&notifications, "created_by = ?", userID)
	if result.Error != nil {
		return []entities.Notification{}, result.Error
	}
	var notificationsEntities []entities.Notification
	for _, notification := range notifications {
		notificationsEntities = append(notificationsEntities, dtos.NotificationModelToEntity(notification))
	}
	return notificationsEntities, nil
}

func (nr *NotificationRepository) DeleteNotificationByID(notificationID uint) (uint, error) {
	var notificationModel dtos.NotificationModel
	result := nr.db.DatabaseConnection.Delete(&notificationModel, notificationID)
	if result.Error != nil {
		return 0, result.Error
	}
	return notificationID, nil
}

func (nr *NotificationRepository) UpdateNotification(notification entities.Notification) (entities.Notification, error) {
	notificationModel := dtos.NotificationToModel(notification)
	result := nr.db.DatabaseConnection.Model(&notificationModel).Where("id = ?", notificationModel.ID).Updates(notificationModel)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	result = nr.db.DatabaseConnection.First(&notificationModel, notificationModel.ID)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	modelEntity := dtos.NotificationModelToEntity(notificationModel)
	return modelEntity, nil
}

func (nr *NotificationRepository) DoesNotificationExistsAndBelongsToUser(userID uint, notificationID uint) (entities.Notification, error) {
	var notificationModel dtos.NotificationModel
	result := nr.db.DatabaseConnection.First(&notificationModel, notificationID)
	if result.Error != nil {
		return entities.Notification{}, errors.New("notification not found")
	}
	if notificationModel.CreatedBy != userID {
		return entities.Notification{}, errors.New("notification does not belong to user")
	}
	notificationEntity := dtos.NotificationModelToEntity(notificationModel)
	return notificationEntity, nil
}

func (nr *NotificationRepository) SetSentAt(notification entities.Notification, time time.Time) (entities.Notification, error) {
	notificationModel := dtos.NotificationToModel(notification)
	result := nr.db.DatabaseConnection.Model(&notificationModel).Where("id = ?", notificationModel.ID).Update("sent_at", time)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	modelEntity := dtos.NotificationModelToEntity(notificationModel)
	return modelEntity, nil
}
