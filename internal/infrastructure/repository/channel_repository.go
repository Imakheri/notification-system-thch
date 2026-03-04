package repository

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
)

type channelRepository struct {
	db *Database
}

func NewChannelRepository(db *Database) gateway.ChannelRepository {
	return &channelRepository{
		db: db,
	}
}

func (cr *channelRepository) GetChannel(channelID uint) (entities.Channel, error) {
	var channelModel dtos.ChannelModel
	result := cr.db.DatabaseConnection.Where("id = ?", channelID).First(&channelModel)
	if result.Error != nil {
		return entities.Channel{}, result.Error
	}
	channelEntity := dtos.ChannelModelToEntity(channelModel)
	return channelEntity, nil
}
