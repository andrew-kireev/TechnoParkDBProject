package forum

import "TechnoParkDBProject/internal/app/forum/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
}
