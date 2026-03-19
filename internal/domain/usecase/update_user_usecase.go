package usecase

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserUseCase interface {
	Exec(userEmail string, user entities.User) (entities.User, error)
}

type updateUserUseCase struct {
	userRepository gateway.UserRepository
}

func NewUpdateUserUseCase(userRepository gateway.UserRepository) UpdateUserUseCase {
	return &updateUserUseCase{
		userRepository: userRepository,
	}
}

func (uuu *updateUserUseCase) Exec(userEmail string, user entities.User) (entities.User, error) {
	userFromDB, err := uuu.userRepository.GetUserByEmail(userEmail)
	if err != nil {
		return entities.User{}, err
	}

	user.ID = userFromDB.ID

	if len(user.Password) > 0 {
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return entities.User{}, errors.New("failed to encrypt password")
		}
		user.Password = string(bytes)
	}

	user, err = uuu.userRepository.UpdateUser(user)
	if err != nil {
		return entities.User{}, err
	}
	updatedUser, err := uuu.userRepository.GetUserByEmail(userEmail)
	if err != nil {
		return entities.User{}, err
	}
	return updatedUser, nil
}
