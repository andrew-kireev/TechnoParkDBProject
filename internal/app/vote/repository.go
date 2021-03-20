package vote

import "TechnoParkDBProject/internal/app/vote/models"

type Repository interface {
	CreateNewVote(vote *models.Vote) error
	UpdateVote(vote *models.Vote) (int, error)
}
