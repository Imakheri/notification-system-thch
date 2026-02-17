package repository

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type userRepository struct {
	db *Database
}

func NewUserRepository(db *Database) gateway.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(user entities.User) (entities.User, error) {
	result := ur.db.DatabaseConnection.Create(&user)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}

func (ur *userRepository) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User
	result := ur.db.DatabaseConnection.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (ur *userRepository) UpdateUser(user entities.User) (entities.User, error) {
	result := ur.db.DatabaseConnection.Save(&user)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}
