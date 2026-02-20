package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	CreatedBy  uint       `gorm:"not null" json:"created_by"`
	SentAt     *time.Time `json:"sent_at"`
	Title      string     `gorm:"type:varchar(100); not null" json:"title"`
	Content    string     `gorm:"not null" json:"content"`
	ChannelID  uint       `json:"channel_id"`
	Channel    Channel    `json:"channel"`
	Recipients []User     `gorm:"many2many:notification_recipients;" json:"recipients"`
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
	channel, err := NewChannel(notification.ChannelID)
	if err != nil {
		return Notification{}, err
	}
	recipients, err := NewRecipients(notification.Recipients)
	if err != nil {
		return Notification{}, err
	}
	notification.Title = title
	notification.Content = content
	notification.ChannelID = channel
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

func NewChannel(channelID uint) (uint, error) {
	if channelID == 0 || channelID > 3 {
		return 0, errors.New("must use a valid channel")
	}
	return channelID, nil
}
