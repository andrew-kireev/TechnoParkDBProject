package repository

import (
	"TechnoParkDBProject/internal/app/posts/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostsRepository struct {
	Conn *pgxpool.Pool
}

func NewPostsRepository(conn *pgxpool.Pool) *PostsRepository {
	return &PostsRepository{
		Conn: conn,
	}
}

func (postRep *PostsRepository) CreatePost(post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (parent, author, message, forum, thread)
			VALUES ($1, $2, $3, $4, $5)
			returning id, parent, author, message, id_edited, forum, thread, created`

	err := postRep.Conn.QueryRow(context.Background(), query, post.Parent, post.Author,
		post.Message, post.Forum, post.Thread).Scan(&post.ID, &post.Parent, &post.Author,
		&post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, err
	}
	return post, nil
}
