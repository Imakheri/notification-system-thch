package usecase

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase interface {
	Exec(user entities.User) (entities.User, error)
}

type createUserUseCase struct {
	repository gateway.UserRepository
}

func NewCreateUser(repository gateway.UserRepository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (cu createUserUseCase) Exec(user entities.User) (entities.User, error) {
	_, err := cu.repository.GetUserByEmail(user.Email)
	if err == nil {
		return entities.User{}, errors.New("user already exists")
	}

	err = checkUserProperties(user)
	if err != nil {
		return entities.User{}, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return entities.User{}, errors.New("failed to encrypt password")
	}
	user.Password = string(bytes)

	user, err = cu.repository.CreateUser(user)
	if err != nil {
		return entities.User{}, err
	}
	user.Password = ""
	return user, nil
}

func isASecurePassword(password string) bool {
	var (
		minLength    bool
		hasUppercase bool
		hasLowercase bool
		hasDigit     bool
		hasSymbol    bool
	)

	if len(password) > 8 && len(password) < 20 {
		minLength = true
	}

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
	return minLength && hasUppercase && hasLowercase && hasDigit && hasSymbol
}

func checkUserProperties(user entities.User) error {
	user.Email = strings.ToLower(user.Email)

	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		return errors.New("invalid email structure")
	}

	if !isASecurePassword(user.Password) {
		return errors.New("invalid password structure")
	}

	if len(user.Name) <= 0 || len(user.Name) >= 20 {
		return errors.New("invalid name structure")
	}

	_, err := strconv.Atoi(user.Phone)

	if len(user.Phone) <= 0 || len(user.Phone) >= 12 || err != nil {
		return errors.New("invalid phone structure")
	}

	if len(user.DeviceToken) <= 0 {
		return errors.New("invalid device token structure")
	}

	return nil
}
