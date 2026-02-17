package gateway

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(user entities.User) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	UpdateUser(user entities.User) (entities.User, error)
}
