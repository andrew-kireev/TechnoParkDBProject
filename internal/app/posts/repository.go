package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Repository interface {
	CreatePost(posts []*models.Post) ([]*models.Post, error)
}
