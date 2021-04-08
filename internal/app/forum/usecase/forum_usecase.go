package usecase

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/forum/models"
	usersModels "TechnoParkDBProject/internal/app/user/models"
)

type ForumUsecase struct {
	forumRep forum.Repository
}

func NewForumUsecase(forumRep forum.Repository) *ForumUsecase {
	return &ForumUsecase{
		forumRep: forumRep,
	}
}

func (forumUse *ForumUsecase) CreateForum(forum *models.Forum) error {
	err := forumUse.forumRep.CreateForum(forum)
	return err
}

func (forumUse *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	forum, err := forumUse.forumRep.GetForumBySlug(slug)
	return forum, err
}

func (forumUse *ForumUsecase) GetUsersByForum(formSlug, since string, limit int, desc bool) ([]*usersModels.User, error) {
	users, err := forumUse.forumRep.GetUsersByForum(formSlug, since, limit, desc)
	if len(users) == 0 {
		_, err := forumUse.forumRep.GetForumBySlug(formSlug)
		if err != nil {
			return nil, err
		}
	}
	return users, err
}
