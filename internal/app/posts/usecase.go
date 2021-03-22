package posts

import "TechnoParkDBProject/internal/app/posts/models"

type Usecase interface {
	CreatePost(posts []*models.Post, slugOrInt string) ([]*models.Post, error)
	GetPosts(sort, since, slugOrID string, limit int, desc bool) ([]*models.Post, error)
	GetPost(postID int, relatedStrs []string) (*models.PostResponse, error)
}
