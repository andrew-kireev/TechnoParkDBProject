package thread

import "TechnoParkDBProject/internal/app/thread/models"

type Usecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	FindThreadBySlug(slug string) (*models.Thread, error)
	GetThreadsByForumSlug(forumSlug, since, desc string, limit int) ([]*models.Thread, error)
}
