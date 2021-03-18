package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Usecase interface {
	CreatePost(posts []*models.Post) ([]*models.Post, error)
}
