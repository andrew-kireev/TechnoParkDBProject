package usecase

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/posts"
	"TechnoParkDBProject/internal/app/posts/models"
	"TechnoParkDBProject/internal/app/thread"
	threadModels "TechnoParkDBProject/internal/app/thread/models"
	"TechnoParkDBProject/internal/app/user"
	"TechnoParkDBProject/internal/pkg/utils"
	"errors"
	"fmt"
	"strconv"
)

type PostsUsecase struct {
	postsRep  posts.Repository
	threadRep thread.Repository
	forumRep  forum.Repository
	userRep   user.Repository
}

func NewPostsUsecase(postsRep posts.Repository, threadRep thread.Repository,
	forumRep forum.Repository, userRep user.Repository) *PostsUsecase {
	return &PostsUsecase{
		postsRep:  postsRep,
		threadRep: threadRep,
		forumRep:  forumRep,
		userRep:   userRep,
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
	if len(posts) == 0 {
		return posts, nil
	}
	if len(posts) != 0 {
		_, err = postUsecase.userRep.GetUserByNickname(posts[0].Author)
		if err != nil {
			return nil, err
		}
	}
	for _, post := range posts {
		post.Thread = thread.ID
		if post.Parent != 0 {
			parentThreadID, err := postUsecase.postsRep.CheckThreadID(post.Parent)
			if err != nil {
				return nil, errors.New("g")
			}
			if parentThreadID != post.Thread {
				return nil, errors.New("g")
			}
		}
		post.Forum = thread.Forum
	}
	posts, err = postUsecase.postsRep.CreatePost(posts)
	return posts, err
}

func (postUsecase *PostsUsecase) GetPosts(sort, since, slugOrID string, limit int, desc bool) ([]*models.Post, error) {
	threadID, err := strconv.Atoi(slugOrID)
	thread := &threadModels.Thread{}
	if sort == "parent_tree" {
		fmt.Println("ura")
	}
	if err != nil {
		thread, err = postUsecase.threadRep.FindThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = postUsecase.threadRep.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}

	posts, err := postUsecase.postsRep.GetPosts(limit, thread.ID, sort, since, desc)
	return posts, err
}

func (posUse *PostsUsecase) GetPost(posID int, relatedStrs []string) (*models.PostResponse, error) {
	post, err := posUse.postsRep.GetPost(posID)
	if err != nil {
		return nil, err
	}
	postResponse := &models.PostResponse{Post: post}
	if utils.Find(relatedStrs, "user") {
		user, err := posUse.userRep.GetUserByNickname(post.Author)
		if err != nil {
			return nil, err
		}
		postResponse.User = user
	}
	if utils.Find(relatedStrs, "forum") {
		forum, err := posUse.forumRep.GetForumBySlug(post.Forum)
		if err != nil {
			return nil, err
		}
		postResponse.Forum = forum
	}
	if utils.Find(relatedStrs, "thread") {
		thread, err := posUse.threadRep.FindThreadByID(post.Thread)
		if err != nil {
			return nil, err
		}
		postResponse.Thread = thread
	}
	return postResponse, err
}

func (postUse *PostsUsecase) UpdatePost(post *models.Post) (*models.Post, error) {
	oldPost, err := postUse.postsRep.GetPost(post.ID)
	if err != nil {
		return nil, err
	}
	if oldPost.Message == post.Message {
		return oldPost, nil
	}
	post, err = postUse.postsRep.UpdatePostByID(post)
	return post, err
}
