package entities

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(50); not null" json:"name"`
	Notifications []Notification `gorm:"many2many:notification_channels;" json:"notifications"`
}
