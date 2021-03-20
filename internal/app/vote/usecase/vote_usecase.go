package usecase

import (
	"TechnoParkDBProject/internal/app/thread"
	threadModels "TechnoParkDBProject/internal/app/thread/models"
	"TechnoParkDBProject/internal/app/vote"
	"TechnoParkDBProject/internal/app/vote/models"
	"strconv"
)

type VoteUsecase struct {
	voteRep   vote.Repository
	threadRep thread.Repository
}

func NewVoteUsecase(voteRep vote.Repository, threadRep thread.Repository) *VoteUsecase {
	return &VoteUsecase{
		voteRep:   voteRep,
		threadRep: threadRep,
	}
}

func (voteUsecase *VoteUsecase) CreateNewVote(vote *models.Vote, slugOrID string) (*threadModels.Thread, error) {
	threadID, err := strconv.Atoi(slugOrID)
	var th *threadModels.Thread
	if err != nil {
		th, err = voteUsecase.threadRep.FindThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	} else {
		th, err = voteUsecase.threadRep.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	vote.ThreadID = th.ID
	err = voteUsecase.voteRep.CreateNewVote(vote)
	th.Votes += 1
	return th, err
}
