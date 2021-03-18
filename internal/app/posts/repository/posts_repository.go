package repository

import (
	"TechnoParkDBProject/internal/app/posts/models"
	"context"
	"fmt"
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

func (postRep *PostsRepository) CreatePost(posts []*models.Post) ([]*models.Post, error) {
	if len(posts) == 0 {
		return posts, nil
	}
	query := `INSERT INTO posts (parent, author, message, forum, thread)
			VALUES `
	for i, post := range posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("(%d, %s, %s, %s, %s)", post.Parent, post.Author,
			post.Message, post.Forum, post.Thread)
	}
	query += "returning id, parent, author, message, is_edited, forum, thread, created"
	fmt.Println(query)

	rows, err := postRep.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}

	return posts, nil
}
