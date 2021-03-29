package posts

import (
	forumModels "TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/posts/models"
)

type Repository interface {
	CreatePost(posts []*models.Post) ([]*models.Post, error)
	FindForumByThreadID(threadID int) (*forumModels.Forum, error)
	GetPosts(limit, threadID int, sort, since string, desc bool) ([]*models.Post, error)
	GetPost(postID int) (*models.Post, error)
	UpdatePostByID(post *models.Post) (*models.Post, error)
	CheckThreadID(parentID int) (int, error)
}
