package repository

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type userRepository struct {
	db_connection *Database
}

func NewUserRepository(db *Database) gateway.UserRepository {
	return &userRepository{
		db_connection: db,
	}
}

func (ur *userRepository) CreateUser(user entities.User) (entities.User, error) {
	result := ur.db_connection.Database().Create(&user)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}

func (ur *userRepository) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User
	result := ur.db_connection.Database().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (ur *userRepository) DoesUserAlreadyExist(email string) bool {
	result := ur.db_connection.Database().Where("email = ?", email).First(&entities.User{})
	if result.Error != nil {
		return false
	}
	return true
}
