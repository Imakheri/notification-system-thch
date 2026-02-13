package usecase

import (
	"errors"
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

	user.Email = strings.ToLower(user.Email)

	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		return entities.User{}, errors.New("invalid email structure")
	}

	if ok := cu.repository.DoesUserAlreadyExist(user.Email); ok {
		return entities.User{}, errors.New("email already exists")
	}

	if !isASecurePassword(user.Password) {
		return entities.User{}, errors.New("invalid password structure")
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

	if len(password) > 8 && len(password) < 16 {
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
