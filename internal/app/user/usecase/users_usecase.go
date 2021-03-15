package usecase

import (
	"TechnoParkDBProject/internal/app/user"
	"TechnoParkDBProject/internal/app/user/models"
)

type UserUsecase struct {
	userRep user.Repository
}

func NewUserUsecase(userRep user.Repository) *UserUsecase {
	return &UserUsecase{
		userRep: userRep,
	}
}

func (userUse *UserUsecase) CreateUser(user *models.User) error {
	err := userUse.userRep.CreateUser(user)
	return err
}

func (userUse *UserUsecase) GetUserByEmailOrNickname(nickname, email string) (*models.User, error) {
	user, err := userUse.userRep.GetUserByEmailOrNickname(nickname, email)
	return user, err
}

func (userUse *UserUsecase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := userUse.userRep.GetUserByNickname(nickname)
	return user, err
}

func (userUse *UserUsecase) UpdateUserInformation(user *models.User) error {
	err := userUse.userRep.UpdateUserInformation(user)
	return err
}

func (userUse *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	user, err := userUse.userRep.GetUserByEmail(email)
	return user, err
}
