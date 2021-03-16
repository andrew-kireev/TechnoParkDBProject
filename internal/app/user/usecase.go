package user

import "TechnoParkDBProject/internal/app/user/models"

type Usecase interface {
	CreateUser(user *models.User) error
	GetUserByEmailOrNickname(nickname, email string) ([]*models.User, error)
	GetUserByNickname(nickname string) (*models.User, error)
	UpdateUserInformation(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}
