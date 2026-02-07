package entities

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	CreatedBy  uint      `gorm:"not null" json:"created_by"`
	Title      string    `gorm:"type:varchar(100); not null" json:"title"`
	Content    string    `gorm:"not null" json:"content"`
	Channels   []Channel `gorm:"many2many:notification_channels;" json:"channels"`
	Recipients []User    `gorm:"many2many:notification_recipients;" json:"recipients"`
}

type NotificationComponent interface {
	Validate() error
	Send() error
}
