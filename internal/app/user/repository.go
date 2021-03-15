package user

import "TechnoParkDBProject/internal/app/user/models"

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByEmailOrNickname(email, nickname string) (*models.User, error)
	GetUserByNickname(nickname string) (*models.User, error)
}
