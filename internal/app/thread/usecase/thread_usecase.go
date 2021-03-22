package usecase

import (
	"TechnoParkDBProject/internal/app/thread"
	"TechnoParkDBProject/internal/app/thread/models"
	"strconv"
)

type ThreadUsecase struct {
	threadRep thread.Repository
}

func NewThreadUsecase(threadRep thread.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		threadRep: threadRep,
	}
}

func (threadUsecase *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	thread, err := threadUsecase.threadRep.CreateThread(thread)
	return thread, err
}

func (threadUsecase *ThreadUsecase) FindThreadBySlug(slug string) (*models.Thread, error) {
	thread, err := threadUsecase.threadRep.FindThreadBySlug(slug)
	return thread, err
}

func (threadUsecase *ThreadUsecase) GetThreadsByForumSlug(forumSlug, since, desc string, limit int) ([]*models.Thread, error) {
	threads, err := threadUsecase.threadRep.GetThreadsByForumSlug(forumSlug, since, desc, limit)
	return threads, err
}

func (threadUsecase *ThreadUsecase) GetThreadBySlugOrID(slugOrID string) (*models.Thread, error) {
	var thread *models.Thread
	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		thread, err = threadUsecase.threadRep.FindThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = threadUsecase.threadRep.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	return thread, err
}

func (threadUsecase *ThreadUsecase) UpdateTreads(slugOrID string, th *models.Thread) (*models.Thread, error) {
	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		th.Slug = slugOrID
		oldThread, errRep := threadUsecase.threadRep.FindThreadBySlug(slugOrID)
		if errRep != nil {
			return nil, errRep
		}
		if th.Title == "" {
			th.Title = oldThread.Title
		}
		if th.Message == "" {
			th.Message = oldThread.Message
		}
		th, errRep = threadUsecase.threadRep.UpdateThreadSlug(th)
		if errRep != nil {
			return nil, errRep
		}
	} else {
		th.ID = threadID
		oldThread, errRep := threadUsecase.threadRep.FindThreadByID(threadID)
		if errRep != nil {
			return nil, errRep
		}
		if th.Title == "" {
			th.Title = oldThread.Title
		}
		if th.Message == "" {
			th.Message = oldThread.Message
		}
		th, errRep = threadUsecase.threadRep.UpdateThreadByID(th)
		if errRep != nil {
			return nil, errRep
		}
	}
	return th, nil
}
