package usecase

import (
	"TechnoParkDBProject/internal/app/posts"
	"TechnoParkDBProject/internal/app/posts/models"
	"TechnoParkDBProject/internal/app/thread"
	threadModels "TechnoParkDBProject/internal/app/thread/models"
	"strconv"
)

type PostsUsecase struct {
	postsRep  posts.Repository
	threadRep thread.Repository
}

func NewPostsUsecase(postsRep posts.Repository, threadRep thread.Repository) *PostsUsecase {
	return &PostsUsecase{
		postsRep:  postsRep,
		threadRep: threadRep,
	}
}

func (postUsecase *PostsUsecase) CreatePost(posts []*models.Post, slugOrInt string) ([]*models.Post, error) {
	threadID, err := strconv.Atoi(slugOrInt)
	thread := &threadModels.Thread{}
	if err != nil {
		thread, err = postUsecase.threadRep.FindThreadBySlug(slugOrInt)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = postUsecase.threadRep.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	for _, post := range posts {
		post.Thread = thread.ID
		post.Forum = thread.Forum
	}
	posts, err = postUsecase.postsRep.CreatePost(posts)
	return posts, err
}
