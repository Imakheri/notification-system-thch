package repository

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
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
	var channel entities.Channel
	result := cr.db.DatabaseConnection.Where("id = ?", channelID).First(&channel)
	if result.Error != nil {
		return entities.Channel{}, result.Error
	}
	return channel, nil
}
