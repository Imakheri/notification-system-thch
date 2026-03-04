package dtos

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Name          string              `gorm:"type:varchar(100); not null" json:"name"`
	Password      string              `gorm:"type:varchar(100); not null" json:"password"`
	Email         string              `gorm:"type:varchar(100); not null" json:"email"`
	Phone         string              `gorm:"type:varchar(100)" json:"phone"`
	DeviceToken   string              `gorm:"type:varchar(100)" json:"device_token"`
	Notifications []NotificationModel `gorm:"many2many:notification_recipients;" json:"notifications"`
}

func (UserModel) TableName() string {
	return "users"
}

func UserModelToEntity(u UserModel) entities.User {
	var notifications []entities.Notification
	for _, notification := range u.Notifications {
		notifications = append(notifications, NotificationModelToEntity(notification))
	}
	return entities.User{
		ID:            u.ID,
		Name:          u.Name,
		Password:      u.Password,
		Email:         u.Email,
		Phone:         u.Phone,
		DeviceToken:   u.DeviceToken,
		Notifications: notifications,
	}
}

func UserToModel(u entities.User) UserModel {
	var notifications []NotificationModel
	for _, notification := range u.Notifications {
		notifications = append(notifications, NotificationToModel(notification))
	}
	return UserModel{
		Model: gorm.Model{
			ID: u.ID,
		},
		Name:          u.Name,
		Password:      u.Password,
		Email:         u.Email,
		Phone:         u.Phone,
		DeviceToken:   u.DeviceToken,
		Notifications: notifications,
	}
}
