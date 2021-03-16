package usecase

import (
	"TechnoParkDBProject/internal/app/forum"
)

type ForumUsecase struct {
	forumRep forum.Repository
}

func NewForumUsecase(forumRep forum.Repository) *ForumUsecase {
	return &ForumUsecase{
		forumRep: forumRep,
	}
}
