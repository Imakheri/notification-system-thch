package usecase

import (
	"errors"

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

	user, err = entities.CheckUserProperties(user)
	if err != nil {
		return entities.User{}, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return entities.User{}, errors.New("failed to encrypt password")
	}
	user.Password = string(bytes)

	newUser, err := cu.repository.CreateUser(user)
	if err != nil {
		return entities.User{}, err
	}
	newUser.Password = ""
	return newUser, nil
}
