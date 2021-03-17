package repository

import (
	"TechnoParkDBProject/internal/app/thread/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ThreadRepository struct {
	Conn *pgxpool.Pool
}

func NewThreadRepository(con *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{
		Conn: con,
	}
}

func (thredRep *ThreadRepository) CreateThread(thread *models.Thread) error {
	query := `INSERT INTO threads (title, author, forum, message, slug)
			VALUES ($1, $2, $3, $4)`

	_, err := thredRep.Conn.Exec(context.Background(), query, thread.Title, thread.Author,
		thread.Forum, thread.Message, thread.Slug)

	return err
}
