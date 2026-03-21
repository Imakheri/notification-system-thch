package dtos

import (
	"time"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type CreateNotificationDTO struct {
	Title      string                          `json:"title"  example:"Example Notification"`
	Content    string                          `json:"content" example:"This is an example notification"`
	ChannelID  uint                            `json:"channel_id" example:"2"`
	Recipients []UserIntoCreateNotificationDTO `json:"recipients"`
}

type NotificationResponseDTO struct {
	ID         uint                      `json:"id" example:"100"`
	CreatedBy  uint                      `json:"created_by" example:"6"`
	SentAt     *time.Time                `json:"sent_at" example:"2026-03-03 16:25:55.470"`
	Title      string                    `json:"title" example:"Example Notification"`
	Content    string                    `json:"content" example:"This is an example notification"`
	ChannelID  uint                      `json:"channel_id" example:"2"`
	Recipients []UserIntoNotificationDTO `json:"recipients"`
}

type UpdateNotificationDTO struct {
	Title      string                          `json:"title" example:"New Example Notification Title"`
	Content    string                          `json:"content" example:"This is an new example notification content"`
	ChannelID  uint                            `json:"channel_id" example:"3"`
	Recipients []UserIntoCreateNotificationDTO `json:"recipients"`
}

type SuccessfulDeleteResponseDTO struct {
	ID      uint   `json:"id" example:"100"`
	Message string `json:"message" example:"Successfully deleted notification"`
}

type ErrorResponseDTO struct {
	Error string `json:"error"`
}

func NotificationToEntity(n CreateNotificationDTO) (entities.Notification, error) {
	var recipients = make([]entities.User, len(n.Recipients))
	for i, recipient := range n.Recipients {
		recipients[i].Email = recipient.Email
	}
	newNotification, err := entities.NewNotification(n.Title, n.Content, n.ChannelID, recipients)
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
	var recipients = make([]UserIntoNotificationDTO, len(n.Recipients))
	for i, recipient := range n.Recipients {
		recipients[i].Email = recipient.Email
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
	var recipients = make([]entities.User, len(n.Recipients))
	for i, recipient := range n.Recipients {
		recipients[i].Email = recipient.Email
	}
	newNotification, err := entities.UpdateNotification(n.Title, n.Content, n.ChannelID, recipients)
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
