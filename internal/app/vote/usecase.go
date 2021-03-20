package vote

import (
	threadModels "TechnoParkDBProject/internal/app/thread/models"
	voteModels "TechnoParkDBProject/internal/app/vote/models"
)

type Usecase interface {
	CreateNewVote(vote *voteModels.Vote, slugOrID string) (*threadModels.Thread, error)
}
