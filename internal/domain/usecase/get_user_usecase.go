package usecase

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/infrastructure/service"
	"golang.org/x/crypto/bcrypt"
)

type GetUser interface {
	Exec(user entities.User) (entities.User, error)
}

type getUser struct {
	repository  gateway.UserRepository
	jwt_service gateway.JwTokenService
}

func NewGetUser(repository gateway.UserRepository) GetUser {
	return &getUser{
		repository: repository,
	}
}

func (g getUser) Exec(userRequest entities.User) (entities.User, error) {
	user, err := g.repository.GetUserByEmail(userRequest.Email)
	if err != nil {
		return entities.User{}, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return entities.User{}, errors.New("entered password is incorrect")
	}

	user.Password = ""

	token, err := service.NewJWTService().GenerateToken(user.Email, user.ID)

	if err != nil {
		return entities.User{}, errors.New("could not generate JWT")
	}

	println(token)
	return user, nil
}
