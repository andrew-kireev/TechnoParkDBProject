package repository

import (
	forumModels "TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/posts/models"
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
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

		query += fmt.Sprintf("(%d, '%s', '%s', '%s', %d)", post.Parent, post.Author,
			post.Message, post.Forum, post.Thread)
	}
	query += "returning id, parent, author, message, is_edited, forum, thread, created"

	rows, err := postRep.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	newPosts := make([]*models.Post, 0)
	defer rows.Close()
	for rows.Next() {
		t := &time.Time{}
		post := &models.Post{}
		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, t)
		post.Created = strfmt.DateTime(t.UTC()).String()
		newPosts = append(newPosts, post)
	}

	return newPosts, nil
}

func (postRep *PostsRepository) FindForumByThreadID(threadID int) (*forumModels.Forum, error) {
	query := `SELECT f.title, f.user_nickname, f.slug FROM forum as f
			left join threads t on f.slug = t.forum
			where t.id = $1`

	forum := &forumModels.Forum{}
	err := postRep.Conn.QueryRow(context.Background(), query,
		threadID).Scan(&forum.Tittle, &forum.UserNickname, &forum.Slug)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (postRep *PostsRepository) GetPosts(limit, threadID int, sort, since string, desc bool) ([]*models.Post, error) {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE thread = $1`
	if since != "" && desc {
		query += " and created <= $2"
	} else if since != "" && !desc {
		query += " and created >= $2"
	} else if since != "" {
		query += " and created >= $2"
	}
	if desc {
		query += " ORDER BY created desc"
	} else if !desc {
		query += " ORDER BY created asc"
	} else {
		query += " ORDER BY created"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)

	fmt.Println(query)

	var rows pgx.Rows
	var err error
	if since != "" {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID, since)
	} else {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		t := &time.Time{}
		post := &models.Post{}
		err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, t)
		post.Created = strfmt.DateTime(t.UTC()).String()
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
