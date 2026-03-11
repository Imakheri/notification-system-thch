package usecase

import (
	"errors"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"golang.org/x/crypto/bcrypt"
)

type GetUserUseCase interface {
	Exec(user entities.User) (entities.User, string, error)
}

type getUserUseCase struct {
	userRepository gateway.UserRepository
	jwtService     gateway.JwTokenService
}

func NewGetUserUseCase(repository gateway.UserRepository, jwtService gateway.JwTokenService) GetUserUseCase {
	return &getUserUseCase{
		userRepository: repository,
		jwtService:     jwtService,
	}
}

func (g *getUserUseCase) Exec(userInput entities.User) (entities.User, string, error) {
	user, err := g.userRepository.GetUserByEmail(userInput.Email)
	if err != nil {
		return entities.User{}, "", errors.New("the e-mail address or password is incorrect")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		return entities.User{}, "", errors.New("the e-mail address or password is incorrect")
	}

	token, err := g.jwtService.GenerateToken(user.Email, user.ID)
	if err != nil {
		return entities.User{}, "", errors.New("could not generate JWT")
	}

	return user, token, nil
}
