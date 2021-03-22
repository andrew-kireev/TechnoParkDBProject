package thread

import "TechnoParkDBProject/internal/app/thread/models"

type Repository interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	FindThreadBySlug(slug string) (*models.Thread, error)
	FindThreadByID(threadID int) (*models.Thread, error)
	GetThreadsByForumSlug(forumSlug, since, desc string, limit int) ([]*models.Thread, error)
	UpdateThreadSlug(thread *models.Thread) (*models.Thread, error)
	UpdateThreadByID(thread *models.Thread) (*models.Thread, error)
}
