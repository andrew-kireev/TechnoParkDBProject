package repository

import (
	"TechnoParkDBProject/internal/app/thread/models"
	"context"
	"fmt"
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

func (thredRep *ThreadRepository) CreateThread(thread *models.Thread) (*models.Thread, error) {
	query := `INSERT INTO threads (title, author, forum, message, slug)
			VALUES ($1, $2, $3, $4, $5)
			Returning id, title, author, forum, message, votes, slug`

	err := thredRep.Conn.QueryRow(context.Background(), query, thread.Title, thread.Author,
		thread.Forum, thread.Message, thread.Slug).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug)

	return thread, err
}

func (thredRep *ThreadRepository) FindThreadBySlug(slug string) (*models.Thread, error) {
	query := `SELECT id, title, author, forum, message, votes, slug, created from threads
			WHERE slug = $1`
	thread := &models.Thread{}

	err := thredRep.Conn.QueryRow(context.Background(), query, slug).Scan(&thread.ID, &thread.Title,
		&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created,
	)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (thredRep *ThreadRepository) GetThreadsByForumSlug(forumSlug, since, desc string) ([]*models.Thread, error) {
	query := `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, created from threads as t
    LEFT JOIN forum f on t.forum = f.slug
	WHERE f.slug = $1 and t.created >= $2`

	rows, err := thredRep.Conn.Query(context.Background(), query, forumSlug, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		err = rows.Scan(&thread.ID, &thread.Title,
			&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			&thread.Slug, &thread.Created)
		if err != nil {
			fmt.Println(err)
		}
		threads = append(threads, thread)
	}
	return threads, nil
}
