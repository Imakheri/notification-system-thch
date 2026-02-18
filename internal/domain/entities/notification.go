package entities

import (
	"errors"

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

type NotificationStrategy interface {
	Send(User, Notification) error
}

func CheckNotificationProperties(notification Notification) (Notification, error) {
	title, err := NewTitle(notification.Title)
	if err != nil {
		return Notification{}, err
	}
	content, err := NewContent(notification.Content)
	if err != nil {
		return Notification{}, err
	}
	channels, err := NewChannels(notification.Channels)
	if err != nil {
		return Notification{}, err
	}
	recipients, err := NewRecipients(notification.Recipients)
	if err != nil {
		return Notification{}, err
	}
	notification.Title = title
	notification.Content = content
	notification.Channels = channels
	notification.Recipients = recipients
	return notification, nil
}

func NewTitle(title string) (string, error) {
	if len(title) == 0 {
		return "", errors.New("notification must have a title")
	}
	return title, nil
}

func NewContent(content string) (string, error) {
	if len(content) == 0 {
		return "", errors.New("notification must have a content")
	}
	return content, nil
}

func NewRecipients(recipients []User) ([]User, error) {
	if len(recipients) < 0 {
		return []User{}, errors.New("recipients must not be greater than 0")
	}
	return recipients, nil
}

func NewChannels(channels []Channel) ([]Channel, error) {
	if len(channels) < 0 {
		return []Channel{}, errors.New("notification must use at least one channel")
	}
	for _, channel := range channels {
		if channel.ID == 0 || channel.ID > 3 {
			return []Channel{}, errors.New("mush enter a valid channel")
		}
	}
	return channels, nil
}
