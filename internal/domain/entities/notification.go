package entities

import (
	"errors"
	"time"
)

type Notification struct {
	ID         uint
	CreatedBy  uint
	SentAt     *time.Time
	Title      string
	Content    string
	ChannelID  uint
	Recipients []User
}

type NotificationStrategy interface {
	Send(string, User, Notification) (int, error)
}

func NewNotification(title string, content string, channelID uint, recipients []User) (Notification, error) {
	validTitle, err := NewTitle(title)
	if err != nil {
		return Notification{}, err
	}
	validContent, err := NewContent(content)
	if err != nil {
		return Notification{}, err
	}
	validChannel, err := NewChannel(channelID)
	if err != nil {
		return Notification{}, err
	}
	validRecipients, err := NewRecipients(recipients)
	if err != nil {
		return Notification{}, err
	}
	return Notification{
		Title:      validTitle,
		Content:    validContent,
		ChannelID:  validChannel,
		Recipients: validRecipients,
	}, nil
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

func UpdateNotification(title string, content string, channelID uint, recipients []User) (Notification, error) {
	validChannel, err := NewChannel(channelID)
	if err != nil {
		return Notification{}, err
	}
	return Notification{
		Title:      title,
		Content:    content,
		ChannelID:  validChannel,
		Recipients: recipients,
	}, nil
}
