package usecase

import (
	"errors"
	"strconv"

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

func (uuu *updateUserUseCase) Exec(userEmail string, userDTO entities.User) (entities.User, error) {
	user, err := uuu.userRepository.GetUserByEmail(userEmail)
	if err != nil {
		return entities.User{}, err
	}

	if len(userDTO.Name) <= 0 {
		userDTO.Name = user.Name
	} else {
		if len(user.Name) >= 20 {
			return entities.User{}, errors.New("invalid name structure")
		}
		user.Name = userDTO.Name
	}

	if len(userDTO.Phone) <= 0 {
		userDTO.Phone = user.Phone
	} else {
		_, err := strconv.Atoi(user.Phone)
		if len(user.Phone) >= 12 || err != nil {
			return entities.User{}, errors.New("invalid phone structure")
		}
		user.Phone = userDTO.Phone
	}

	if len(userDTO.DeviceToken) <= 0 {
		userDTO.DeviceToken = user.DeviceToken
	} else {
		user.DeviceToken = userDTO.DeviceToken
	}

	if len(userDTO.Password) <= 0 {
		userDTO.Password = user.Password
	} else {
		newPassword, err := entities.NewPassword(userDTO.Password)
		if err != nil {
			return entities.User{}, err
		}
		bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
		if err != nil {
			return entities.User{}, errors.New("failed to encrypt password")
		}
		user.Password = string(bytes)
	}

	user, err = uuu.userRepository.UpdateUser(user)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
