package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type User struct {
	ID            uint
	Name          string
	Password      string
	Email         string
	Phone         string
	DeviceToken   string
	Notifications []Notification
}

func NewUser(name, password, email, phone, deviceToken string) (User, error) {
	validName, err := NewName(name)
	if err != nil {
		return User{}, err
	}
	validPassword, err := NewPassword(password)
	if err != nil {
		return User{}, err
	}
	validEmail, err := NewEmail(email)
	if err != nil {
		return User{}, err
	}
	validPhone, err := NewPhone(phone)
	if err != nil {
		return User{}, err
	}
	validDeviceToken, err := NewDeviceToken(deviceToken)
	if err != nil {
		return User{}, err
	}
	return User{
		Name:        validName,
		Password:    validPassword,
		Email:       validEmail,
		Phone:       validPhone,
		DeviceToken: validDeviceToken,
	}, nil
}

func NewName(name string) (string, error) {
	if len(name) <= 0 || len(name) >= 25 {
		return "", errors.New("invalid name structure")
	}
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsSpace(ch) {
			return "", errors.New("name should only contain letters")
		}
	}
	return name, nil
}

func NewEmail(email string) (string, error) {
	email = strings.ToLower(email)

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") || len(email) < 10 {
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

func UpdateUser(name, password, phone, deviceToken string) (User, error) {
	validName, err := UpdateName(name)
	if err != nil {
		return User{}, err
	}
	validPassword, err := UpdatePassword(password)
	if err != nil {
		return User{}, err
	}
	validPhone, err := UpdatePhone(phone)
	if err != nil {
		return User{}, err
	}
	return User{
		Name:        validName,
		Password:    validPassword,
		Phone:       validPhone,
		DeviceToken: deviceToken,
	}, nil
}

func UpdateName(name string) (string, error) {
	if len(name) >= 25 {
		return "", errors.New("invalid name structure")
	}
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsSpace(ch) {
			return "", errors.New("name should only contain letters")
		}
	}
	return name, nil
}

func UpdatePassword(password string) (string, error) {
	if len(password) == 0 {
		return "", nil
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

func UpdatePhone(phone string) (string, error) {
	if len(phone) >= 12 {
		return "", errors.New("invalid phone structure")
	}
	return phone, nil
}
