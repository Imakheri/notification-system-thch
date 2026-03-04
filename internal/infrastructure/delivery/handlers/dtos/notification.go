package dtos

import (
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type CreateNotificationDTO struct {
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	ChannelID  uint            `json:"channel_id"`
	Recipients []entities.User `json:"recipients"`
}

type NotificationResponseDTO struct {
	ID         uint                      `json:"id"`
	CreatedBy  uint                      `json:"created_by"`
	SentAt     *time.Time                `json:"sent_at"`
	Title      string                    `json:"title"`
	Content    string                    `json:"content"`
	ChannelID  uint                      `json:"channel_id"`
	Recipients []UserIntoNotificationDTO `json:"recipients"`
}

type UpdateNotificationDTO struct {
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	ChannelID  uint            `json:"channel_id"`
	Recipients []entities.User `json:"recipients"`
}

func NotificationToEntity(n CreateNotificationDTO) (entities.Notification, error) {
	newNotification, err := entities.NewNotification(n.Title, n.Content, n.ChannelID, n.Recipients)
	if err != nil {
		return entities.Notification{}, err
	}
	newNotificationEntity := entities.Notification{
		Title:      newNotification.Title,
		Content:    newNotification.Content,
		ChannelID:  newNotification.ChannelID,
		Recipients: newNotification.Recipients,
	}
	return newNotificationEntity, nil
}

func NotificationToDto(n entities.Notification) NotificationResponseDTO {
	var recipients []UserIntoNotificationDTO

	for _, recipient := range n.Recipients {
		recipients = append(recipients, UserIntoNotificationToDTO(recipient))
	}
	return NotificationResponseDTO{
		ID:         n.ID,
		CreatedBy:  n.CreatedBy,
		Title:      n.Title,
		Content:    n.Content,
		ChannelID:  n.ChannelID,
		Recipients: recipients,
		SentAt:     n.SentAt,
	}
}

func NotificationUpdateToEntity(n UpdateNotificationDTO) (entities.Notification, error) {
	newNotification, err := entities.UpdateNotification(n.Title, n.Content, n.ChannelID, n.Recipients)
	if err != nil {
		return entities.Notification{}, err
	}
	newNotificationEntity := entities.Notification{
		Title:      newNotification.Title,
		Content:    newNotification.Content,
		ChannelID:  newNotification.ChannelID,
		Recipients: newNotification.Recipients,
	}
	return newNotificationEntity, nil
}
