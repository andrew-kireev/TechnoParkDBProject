package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Usecase interface {
	CreatePost(posts []*models.Post, slugOrInt string) ([]*models.Post, error)
}
