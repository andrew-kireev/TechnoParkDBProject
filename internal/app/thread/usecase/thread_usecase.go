package usecase

import (
	"TechnoParkDBProject/internal/app/thread"
	"TechnoParkDBProject/internal/app/thread/models"
)

type ThreadUsecase struct {
	threadRep thread.Repository
}

func NewThreadUsecase(threadRep thread.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		threadRep: threadRep,
	}
}

func (threadUsecase *ThreadUsecase) CreateThread(thread *models.Thread) error {
	err := threadUsecase.threadRep.CreateThread(thread)
	return err
}
