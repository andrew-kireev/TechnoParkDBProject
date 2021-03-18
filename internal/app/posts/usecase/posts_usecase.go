package usecase

import (
	"TechnoParkDBProject/internal/app/posts"
	"TechnoParkDBProject/internal/app/posts/models"
)

type PostsUsecase struct {
	postsRep posts.Repository
}

func NewPostsUsecase(postsRep posts.Repository) *PostsUsecase {
	return &PostsUsecase{
		postsRep: postsRep,
	}
}

func (postUsecase *PostsUsecase) CreatePost(posts []*models.Post) ([]*models.Post, error) {
	pos, err := postUsecase.postsRep.CreatePost(posts)
	return pos, err
}
