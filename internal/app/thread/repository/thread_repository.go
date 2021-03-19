package repository

import (
	"TechnoParkDBProject/internal/app/thread/models"
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
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
	var err error
	if thread.Created != "" {
		query := `INSERT INTO threads (title, author, forum, message, slug, created)
			VALUES ($1, $2, $3, $4, $5, $6)
			Returning id, title, author, forum, message, votes, slug`

		err = thredRep.Conn.QueryRow(context.Background(), query, thread.Title, thread.Author,
			thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(&thread.ID, &thread.Title, &thread.Author,
			&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug)
	} else {
		query := `INSERT INTO threads (title, author, forum, message, slug)
			VALUES ($1, $2, $3, $4, $5)
			Returning id, title, author, forum, message, votes, slug`

		err = thredRep.Conn.QueryRow(context.Background(), query, thread.Title, thread.Author,
			thread.Forum, thread.Message, thread.Slug).Scan(&thread.ID, &thread.Title, &thread.Author,
			&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug)
	}
	return thread, err
}

func (thredRep *ThreadRepository) FindThreadBySlug(slug string) (*models.Thread, error) {
	query := `SELECT id, title, author, forum, message, votes, slug, created from threads
			WHERE slug = $1`
	thread := &models.Thread{}
	t := &time.Time{}

	err := thredRep.Conn.QueryRow(context.Background(), query, slug).Scan(&thread.ID, &thread.Title,
		&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, t,
	)
	thread.Created = strfmt.DateTime(t.UTC()).String()
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (thredRep *ThreadRepository) FindThreadByID(threadID int) (*models.Thread, error) {
	query := `SELECT id, title, author, forum, message, votes, slug, created from threads
			WHERE id = $1`
	thread := &models.Thread{}
	t := &time.Time{}

	err := thredRep.Conn.QueryRow(context.Background(), query, threadID).Scan(&thread.ID, &thread.Title,
		&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, t,
	)
	thread.Created = strfmt.DateTime(t.UTC()).String()
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (thredRep *ThreadRepository) GetThreadsByForumSlug(forumSlug, since, desc string, limit int) ([]*models.Thread, error) {
	query := `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t
  LEFT JOIN forum f on t.forum = f.slug
	WHERE f.slug = $1`
	if since != "" && desc == "true" {
		query += " and t.created <= $2"
	} else if since != "" && desc == "false" {
		query += " and t.created >= $2"
	} else if since != "" {
		query += " and t.created >= $2"
	}
	if desc == "true" {
		query += " ORDER BY t.created desc"
	} else if desc == "false" {
		query += " ORDER BY t.created asc"
	} else {
		query += " ORDER BY t.created"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	var rows pgx.Rows
	var err error
	if since != "" {
		rows, err = thredRep.Conn.Query(context.Background(), query, forumSlug, since)
	} else {
		rows, err = thredRep.Conn.Query(context.Background(), query, forumSlug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		t := &time.Time{}
		thread := &models.Thread{}
		err = rows.Scan(&thread.ID, &thread.Title,
			&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			&thread.Slug, t)
		thread.Created = strfmt.DateTime(t.UTC()).String()
		//if err != nil {
		//	fmt.Println(err)
		//}
		threads = append(threads, thread)
	}
	return threads, nil
}
