package entities

import (
	"errors"
	"time"
	"unicode"
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
	if len(title) < 10 {
		return "", errors.New("title notification must have at least 10 characters")
	}
	var digits int
	var symbols int
	for _, ch := range title {
		if unicode.IsDigit(ch) {
			digits++
		}
		if unicode.IsSymbol(ch) {
			symbols++
		}
		if digits > 6 || symbols > 4 {
			return "", errors.New("title notification must not have more than 6 digits or symbols")
		}
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
	if recipients == nil {
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
	validTitle, err := UpdateTitle(title)
	if err != nil {
		return Notification{}, err
	}
	validChannel, err := NewChannel(channelID)
	if err != nil {
		return Notification{}, err
	}
	return Notification{
		Title:      validTitle,
		Content:    content,
		ChannelID:  validChannel,
		Recipients: recipients,
	}, nil
}

func UpdateTitle(title string) (string, error) {
	if len(title) == 0 {
		return "", nil
	}
	if len(title) < 10 {
		return "", errors.New("title notification must have at least 10 characters")
	}
	var digits int
	var symbols int
	for _, ch := range title {
		if unicode.IsDigit(ch) {
			digits++
		}
		if unicode.IsSymbol(ch) {
			symbols++
		}
		if digits > 6 || symbols > 4 {
			return "", errors.New("title notification must not have more than 6 digits or symbols")
		}
	}
	return title, nil
}
