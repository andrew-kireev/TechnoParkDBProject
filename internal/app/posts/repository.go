package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Repository interface {
	CreatePost(post *models.Post) (*models.Post, error)
}
