package usecase

import (
	"TechnoParkDBProject/internal/app/thread"
	threadModels "TechnoParkDBProject/internal/app/thread/models"
	"TechnoParkDBProject/internal/app/vote"
	"TechnoParkDBProject/internal/app/vote/models"
	"github.com/jackc/pgconn"
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
	if err != nil {
		if err.(*pgconn.PgError).Code == "23503" {
			return nil, err
		}
		linesUpdated, err := voteUsecase.voteRep.UpdateVote(vote)
		if err != nil {
			return nil, err
		}
		if linesUpdated != 0 {
			th.Votes += 2 * vote.Voice
		}
		return th, err
	}
	th.Votes += vote.Voice
	return th, err
}
