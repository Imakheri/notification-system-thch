package usecase

import (
	"errors"
	"log"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/infrastructure/strategies"
)

type CreateNotificationUseCase interface {
	Exec(userID uint, userEmail string, notification entities.Notification) (entities.Notification, error)
}

type createNotificationUseCase struct {
	userRepository         gateway.UserRepository
	notificationRepository gateway.NotificationRepository
	channelRepository      gateway.ChannelRepository
	simulatedApiService    gateway.SimulatedApiService
	clock                  gateway.Clock
}

func NewCreateNotificationUseCase(notificationRepository gateway.NotificationRepository, userRepository gateway.UserRepository, channelRepository gateway.ChannelRepository, simulatedApiService gateway.SimulatedApiService, clock gateway.Clock) CreateNotificationUseCase {
	return &createNotificationUseCase{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		channelRepository:      channelRepository,
		simulatedApiService:    simulatedApiService,
		clock:                  clock,
	}
}

func (cnu *createNotificationUseCase) Exec(userID uint, userEmail string, notification entities.Notification) (entities.Notification, error) {
	channel, err := cnu.channelRepository.GetChannel(notification.ChannelID)
	notification.ChannelID = channel.ID

	for i, recipient := range notification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			return entities.Notification{}, errors.New("recipient does not exist")
		}
		if userID == user.ID {
			return entities.Notification{}, errors.New("invalid recipient")
		}
		notification.Recipients[i].ID = user.ID
	}

	sender, err := cnu.userRepository.GetUserByEmail(userEmail)
	if err != nil {
		return entities.Notification{}, errors.New("cannot get sender information")
	}

	newNotification, err := cnu.notificationRepository.CreateNotification(userID, notification)
	if err != nil {
		return entities.Notification{}, err
	}
	for _, recipient := range newNotification.Recipients {
		user, err := cnu.userRepository.GetUserByEmail(recipient.Email)
		if err != nil {
			_, err := cnu.notificationRepository.DeleteNotificationByID(newNotification.ID)
			return entities.Notification{}, err
		}
		currentStrategy := cnu.strategySelector(notification.ChannelID)
		_, err = currentStrategy.Send(sender.Email, user, newNotification)
		if err != nil {
			_, errorDeleting := cnu.notificationRepository.DeleteNotificationByID(newNotification.ID)
			if errorDeleting != nil {
				return entities.Notification{}, errorDeleting
			}
			return entities.Notification{}, err
		}
	}
	newNotification, err = cnu.notificationRepository.SetSentAt(newNotification, cnu.clock.Now())
	if err != nil {
		return entities.Notification{}, err
	}

	return newNotification, nil
}

func (cnu *createNotificationUseCase) strategySelector(idStrategy uint) entities.NotificationStrategy {
	var selectedStrategy entities.NotificationStrategy
	switch idStrategy {
	case 1:
		selectedStrategy = strategies.NewEmailStrategy(cnu.simulatedApiService)
	case 2:
		selectedStrategy = strategies.NewSMSStrategy(cnu.simulatedApiService)
	case 3:
		selectedStrategy = strategies.NewPushStrategy(cnu.simulatedApiService)
	default:
		log.Fatal("Unknown strategy")
	}
	return selectedStrategy
}
