package repository

import (
	"TechnoParkDBProject/internal/app/vote/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VoteRepository struct {
	Conn *pgxpool.Pool
}

func NewVoteRepository(con *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{
		Conn: con,
	}
}

func (voteRep *VoteRepository) CreateNewVote(vote *models.Vote) error {
	query := `INSERT INTO votes (nickname, thread_id, voice)
			VALUES ($1, $2, $3)`

	_, err := voteRep.Conn.Exec(context.Background(), query, vote.Nickname,
		vote.ThreadID, vote.Voice)
	return err
}

func (voteRep *VoteRepository) UpdateVote(vote *models.Vote) (int, error) {
	query := `UPDATE votes SET voice = $1
		WHERE thread_id = $2 and nickname = $3 and voice != $1`

	res, err := voteRep.Conn.Exec(context.Background(), query, vote.Voice, vote.ThreadID, vote.Nickname)
	return int(res.RowsAffected()), err
}
