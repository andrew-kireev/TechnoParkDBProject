package usecase

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/forum/models"
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
