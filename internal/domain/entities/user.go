package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(100); not null" json:"name"`
	Password      string         `gorm:"type:varchar(100); not null" json:"password"`
	Email         string         `gorm:"type:varchar(100); not null" json:"email"`
	Phone         string         `gorm:"type:varchar(100)" json:"phone"`
	DeviceToken   string         `gorm:"type:varchar(100)" json:"device_token"`
	Notifications []Notification `gorm:"many2many:notification_recipients;" json:"notifications"`
}

func CheckUserProperties(user User) (User, error) {
	name, err := NewName(user.Name)
	if err != nil {
		return User{}, err
	}
	password, err := NewPassword(user.Password)
	if err != nil {
		return User{}, err
	}
	email, err := NewEmail(user.Email)
	if err != nil {
		return User{}, err
	}
	phone, err := NewPhone(user.Phone)
	if err != nil {
		return User{}, err
	}
	deviceToken, err := NewDeviceToken(user.DeviceToken)
	if err != nil {
		return User{}, err
	}
	user.Name = name
	user.Password = password
	user.Email = email
	user.Phone = phone
	user.DeviceToken = deviceToken
	return user, nil
}

func NewName(name string) (string, error) {
	if len(name) <= 0 || len(name) >= 20 {
		return "", errors.New("invalid name structure")
	}
	return name, nil
}

func NewEmail(email string) (string, error) {
	email = strings.ToLower(email)

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return "", errors.New("invalid email structure")
	}

	return email, nil
}

func NewPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	if len(password) > 20 {
		return "", errors.New("password must be less than 20 characters")
	}

	var (
		hasUppercase bool
		hasLowercase bool
		hasDigit     bool
		hasSymbol    bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUppercase || !hasLowercase || !hasDigit || !hasSymbol {
		return "", errors.New(fmt.Sprintf("password must contain at least one upper case and one lower case character, one digit and one special character"))
	}

	return password, nil
}

func NewPhone(phone string) (string, error) {
	_, err := strconv.Atoi(phone)
	if len(phone) <= 0 || len(phone) >= 12 || err != nil {
		return "", errors.New("invalid phone structure")
	}
	return phone, nil
}

func NewDeviceToken(deviceToken string) (string, error) {
	if len(deviceToken) <= 0 {
		return "", errors.New("invalid device token structure")
	}
	return deviceToken, nil
}
