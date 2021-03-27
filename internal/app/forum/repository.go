package forum

import (
	"TechnoParkDBProject/internal/app/forum/models"
	usersModels "TechnoParkDBProject/internal/app/user/models"
)

type Repository interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
	GetUsersByForum(formSlug, since string, limit int, desc bool) ([]*usersModels.User, error)
}
