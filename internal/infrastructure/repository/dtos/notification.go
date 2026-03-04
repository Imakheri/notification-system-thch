package dtos

import (
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"gorm.io/gorm"
)

type NotificationModel struct {
	gorm.Model
	CreatedBy  uint         `gorm:"not null" json:"created_by"`
	SentAt     *time.Time   `json:"sent_at"`
	Title      string       `gorm:"type:varchar(100); not null" json:"title"`
	Content    string       `gorm:"not null" json:"content"`
	ChannelID  uint         `json:"channel_id"`
	Channel    ChannelModel `json:"channel"`
	Recipients []UserModel  `gorm:"many2many:notification_recipients;" json:"recipients"`
}

func (NotificationModel) TableName() string {
	return "notifications"
}

func NotificationModelToEntity(n NotificationModel) entities.Notification {
	var recipients []entities.User
	for _, recipient := range n.Recipients {
		recipients = append(recipients, UserModelToEntity(recipient))
	}
	return entities.Notification{
		ID:         n.ID,
		CreatedBy:  n.CreatedBy,
		SentAt:     n.SentAt,
		Title:      n.Title,
		Content:    n.Content,
		ChannelID:  n.ChannelID,
		Recipients: recipients,
	}
}

func NotificationToModel(n entities.Notification) NotificationModel {
	var recipients []UserModel
	for _, recipient := range n.Recipients {
		recipients = append(recipients, UserToModel(recipient))
	}
	return NotificationModel{
		Model: gorm.Model{
			ID: n.ID,
		},
		CreatedBy:  n.CreatedBy,
		SentAt:     n.SentAt,
		Title:      n.Title,
		Content:    n.Content,
		ChannelID:  n.ChannelID,
		Recipients: recipients,
	}
}
