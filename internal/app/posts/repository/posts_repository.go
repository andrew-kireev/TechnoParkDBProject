package repository

import (
	forumModels "TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/posts/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
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
	query := `INSERT INTO posts (parent, author, message, forum, thread)
			VALUES `
	for i, post := range posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("(NULLIF(%d, 0), '%s', '%s', '%s', %d)", post.Parent, post.Author,
			post.Message, post.Forum, post.Thread)
	}
	query += " returning id, parent, author, message, is_edited, forum, thread, created"
	rows, err := postRep.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	newPosts := make([]*models.Post, 0)
	var parent sql.NullInt64
	defer rows.Close()
	for rows.Next() {
		t := &time.Time{}
		post := &models.Post{}
		err = rows.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, t)
		if err != nil {
			fmt.Println(err)
		}
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			return nil, err
		}
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
	postID, _ := strconv.Atoi(since)
	var query string

	if sort == "flat" || sort == "" {
		query = FormQueryFlatSort(limit, threadID, sort, since, desc)
	} else if sort == "tree" {
		query = FormQuerySortTree(limit, threadID, postID, sort, since, desc)
	} else if sort == "parent_tree" {
		query = FormQuerySortParentTree(limit, threadID, postID, sort, since, desc)
	}
	if sort != "parent_tree" {
		query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	}
	var rows pgx.Rows
	var err error
	if sort == "tree" {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID)
	} else if sort == "parent_tree" {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID)
	} else if since != "" {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID, postID)
	} else {
		rows, err = postRep.Conn.Query(context.Background(), query, threadID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	var parent sql.NullInt64
	for rows.Next() {
		t := &time.Time{}
		post := &models.Post{}
		err = rows.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, t)
		post.Created = strfmt.DateTime(t.UTC()).String()
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func FormQueryFlatSort(limit, threadID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE thread = $1`
	if since != "" && desc {
		query += " and id < $2"
	} else if since != "" && !desc {
		query += " and id > $2"
	} else if since != "" {
		query += " and id > $2"
	}
	if desc {
		query += " ORDER BY created desc, id desc"
	} else if !desc {
		query += " ORDER BY created asc, id"
	} else {
		query += " ORDER BY created, id"
	}
	return query
}

func FormQuerySortTree(limit, threadID, ID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE thread = $1`
	if since != "" && desc {
		query += " and path < "
	} else if since != "" && !desc {
		query += " and path > "
	} else if since != "" {
		query += " and path > "
	}
	if since != "" {
		query += fmt.Sprintf(` (SELECT path FROM posts WHERE id = %d) `, ID)
	}
	if desc {
		query += " ORDER BY path desc"
	} else if !desc {
		query += " ORDER BY path asc, id"
	} else {
		query += " ORDER BY path, id"
	}
	return query
}

func FormQuerySortParentTree(limit, threadID, ID int, sort, since string, desc bool) string {
	subQuery := ""
	query := `select id, parent, author, message, is_edited, forum, thread,
				created FROM posts WHERE path[1] IN `
	if since != "" {
		if desc {
			subQuery = `and path[1] < `
		} else {
			subQuery = `and path[1] > `
		}
		subQuery += fmt.Sprintf(`(select path[1] from posts where id = %d)`, ID)
	}
	subQuery = `select id FROM posts WHERE thread = $1 and parent is null ` + subQuery
	if desc {
		subQuery += `order by id desc`
		subQuery += fmt.Sprintf(` LIMIT NULLIF(%d, 0)`, limit)
		query += fmt.Sprintf(`(%s) ORDER BY path[1] DESC, path, id`, subQuery)
	} else {
		subQuery += `order by id asc`
		subQuery += fmt.Sprintf(` LIMIT NULLIF(%d, 0)`, limit)
		query += fmt.Sprintf(`(%s) order by path, id`, subQuery)
	}
	return query
}

func (postRep *PostsRepository) GetPost(postID int) (*models.Post, error) {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE id = $1
			order by created, id`

	post := &models.Post{}
	t := &time.Time{}
	var parent sql.NullInt64
	err := postRep.Conn.QueryRow(context.Background(), query, postID).Scan(&post.ID,
		&parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum,
		&post.Thread, t)
	post.Created = strfmt.DateTime(t.UTC()).String()
	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	if err != nil {
		return nil, err
	}
	return post, nil
}

func (postRep *PostsRepository) UpdatePostByID(post *models.Post) (*models.Post, error) {
	query := `UPDATE posts
SET message   = (CASE
                     WHEN LTRIM($1) = '' THEN message
                     ELSE $1 END),
    is_edited = (CASE
                     WHEN LTRIM($1) = '' THEN false
                     ELSE true END)
WHERE id = $2
returning id, parent, author, message, is_edited, forum, thread, created `

	t := &time.Time{}
	var parent sql.NullInt64
	err := postRep.Conn.QueryRow(context.Background(), query, post.Message, post.ID).Scan(
		&post.ID, &parent, &post.Author, &post.Message,
		&post.IsEdited, &post.Forum, &post.Thread, t)
	if err != nil {
		return nil, err
	}
	post.Created = strfmt.DateTime(t.UTC()).String()
	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	return post, nil
}

func (postRep *PostsRepository) CheckThreadID(parentID int) (int, error) {
	query := `SELECT thread from posts where id = $1`
	var threadID int

	err := postRep.Conn.QueryRow(context.Background(), query, parentID).Scan(&threadID)
	return threadID, err
}
