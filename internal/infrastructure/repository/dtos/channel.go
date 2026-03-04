package dtos

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"gorm.io/gorm"
)

type ChannelModel struct {
	gorm.Model
	Name string `gorm:"type:varchar(50); not null" json:"name"`
}

func (ChannelModel) TableName() string {
	return "channels"
}

func ChannelModelToEntity(ch ChannelModel) entities.Channel {
	return entities.Channel{
		ID:   ch.ID,
		Name: ch.Name,
	}
}
