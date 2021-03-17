package thread

import "TechnoParkDBProject/internal/app/thread/models"

type Repository interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	FindThreadBySlug(slug string) (*models.Thread, error)
	GetThreadsByForumSlug(forumSlug, since, desc string) ([]*models.Thread, error)
}
