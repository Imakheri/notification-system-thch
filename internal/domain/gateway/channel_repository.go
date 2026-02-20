package gateway

import "github.com/imakheri/notifications-thch/internal/domain/entities"

type ChannelRepository interface {
	GetChannel(uint) (entities.Channel, error)
}
