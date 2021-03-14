package user

import "TechnoParkDBProject/internal/app/user/models"

type Useacse interface {
	CreateUser(user *models.User) error
	GetUserByEmailOrNickname(nickname, email string) (*models.User, error)
}