package repository

import (
	"TechnoParkDBProject/internal/app/forum/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ForumRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(con *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		Conn: con,
	}
}

func (forumRep *ForumRepository) CreateForum(forum *models.Forum) error {
	query := `INSERT INTO forum (title, user_nickname, slug, posts, threads)
			VALUES ($1, $2, $3, $4, $5)`

	_, err := forumRep.Conn.Exec(context.Background(), query, forum.Tittle,
		forum.UserNickname, forum.Slug, forum.Posts, forum.Threads,
	)
	return err
}

func (forumRep *ForumRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads from forum
			where slug = $1`
	forum := &models.Forum{}

	err := forumRep.Conn.QueryRow(context.Background(), query,
		slug).Scan(&forum.Tittle, &forum.UserNickname, &forum.Slug,
		&forum.Posts, &forum.Threads)

	if err != nil {
		return nil, err
	}
	return forum, nil
}
