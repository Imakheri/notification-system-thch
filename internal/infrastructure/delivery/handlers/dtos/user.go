package dtos

import "github.com/imakheri/notifications-thch/internal/domain/entities"

type CreateUserDTO struct {
	Name        string `json:"name" example:"William Shakespeare"`
	Password    string `json:"password" example:"Hamlet123."`
	Email       string `json:"email" example:"williams@example.com"`
	Phone       string `json:"phone" example:"1234567890"`
	DeviceToken string `json:"device_token" example:"0q1e2r3t4y5a6b7c8d9e"`
}

type ResponseUserDTO struct {
	Name        string `json:"name" example:"William Shakespeare"`
	Email       string `json:"email" example:"williams@example.com"`
	Phone       string `json:"phone" example:"1234567890"`
	DeviceToken string `json:"device_token" example:"0q1e2r3t4y5a6b7c8d9e"`
}

type UpdateUserDTO struct {
	Name        string `json:"name" example:"Fernando Pessoa"`
	Password    string `json:"password" example:"Mensagem123."`
	Phone       string `json:"phone" example:"0129834765"`
	DeviceToken string `json:"device_token" example:"12345qwerty67890"`
}

type UserIntoCreateNotificationDTO struct {
	Email string `json:"email" example:"williams@example.com"`
}

type UserIntoNotificationDTO struct {
	ID    uint   `json:"id" example:"20"`
	Email string `json:"email" example:"williams@example.com"`
}

type GetUserDTO struct {
	Email    string `json:"email" example:"williams@example.com"`
	Password string `json:"password" example:"Hamlet123."`
}

type LoginResponseDTO struct {
	User  ResponseUserDTO `json:"user"`
	Token string          `json:"token" `
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

func GetUserToEntity(u GetUserDTO) entities.User {
	return entities.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
