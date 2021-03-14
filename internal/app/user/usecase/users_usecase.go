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

func (UserUse *UserUsecase) CreateUser(user *models.User) error {
	err := UserUse.userRep.CreateUser(user)
	return err
}

func (UserUse *UserUsecase) GetUserByEmailOrNickname(nickname, email string) (*models.User, error) {
	user, err := UserUse.userRep.GetUserByEmailOrNickname(nickname, email)
	return user, err
}
