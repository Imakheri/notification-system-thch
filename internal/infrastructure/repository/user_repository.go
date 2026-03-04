package repository

import (
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
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
	userModel := dtos.UserToModel(user)
	result := ur.db.DatabaseConnection.Create(&userModel)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	userEntity := dtos.UserModelToEntity(userModel)
	return userEntity, nil
}

func (ur *userRepository) GetUserByEmail(email string) (entities.User, error) {
	var userModel dtos.UserModel
	result := ur.db.DatabaseConnection.Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	userEntity := dtos.UserModelToEntity(userModel)
	return userEntity, nil
}

func (ur *userRepository) UpdateUser(user entities.User) (entities.User, error) {
	userModel := dtos.UserToModel(user)
	result := ur.db.DatabaseConnection.Model(&userModel).Where("id = ?", userModel.ID).Updates(userModel)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	userEntity := dtos.UserModelToEntity(userModel)
	return userEntity, nil
}
