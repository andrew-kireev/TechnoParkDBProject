package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Usecase interface {
	CreatePost(post *models.Post) (*models.Post, error)
}
