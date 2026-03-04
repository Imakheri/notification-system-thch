package dtos

import "github.com/imakheri/notifications-thch/internal/domain/entities"

type CreateUserDTO struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DeviceToken string `json:"device_token"`
}

type ResponseUserDTO struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DeviceToken string `json:"device_token"`
}

type UpdateUserDTO struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DeviceToken string `json:"device_token"`
}

type UserIntoNotificationDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type GetUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserToEntity(u CreateUserDTO) (entities.User, error) {
	newUser, err := entities.NewUser(u.Name, u.Password, u.Email, u.Phone, u.DeviceToken)
	if err != nil {
		return entities.User{}, err
	}
	newUserEntity := entities.User{
		Name:        newUser.Name,
		Password:    newUser.Password,
		Email:       newUser.Email,
		Phone:       newUser.Phone,
		DeviceToken: newUser.DeviceToken,
	}
	return newUserEntity, nil
}

func UserUpdateToEntity(u UpdateUserDTO) (entities.User, error) {
	updateUser, err := entities.UpdateUser(u.Name, u.Password, u.Phone, u.DeviceToken)
	if err != nil {
		return entities.User{}, err
	}
	updatedUserEntity := entities.User{
		Name:        updateUser.Name,
		Password:    updateUser.Password,
		Phone:       updateUser.Phone,
		DeviceToken: updateUser.DeviceToken,
	}
	return updatedUserEntity, nil
}

func UserResponseToDTO(u entities.User) ResponseUserDTO {
	return ResponseUserDTO{
		Name:        u.Name,
		Email:       u.Email,
		Phone:       u.Phone,
		DeviceToken: u.DeviceToken,
	}
}

func UserIntoNotificationToDTO(u entities.User) UserIntoNotificationDTO {
	return UserIntoNotificationDTO{
		ID:    u.ID,
		Email: u.Email,
	}
}

func GetUserToEntity(u GetUserDTO) entities.User {
	return entities.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
