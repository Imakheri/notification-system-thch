package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(100); not null" json:"name"`
	Password      string         `gorm:"type:varchar(100); not null" json:"password"`
	Email         string         `gorm:"type:varchar(100); not null" json:"email"`
	Notifications []Notification `gorm:"many2many:notification_recipients;" json:"notifications"`
}
