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

func (postRep *PostsRepository) CreatePost(posts []*models.Post, forum string, threadID int) ([]*models.Post, error) {
	if posts[0].Parent != 0 {
		pThread, err := postRep.CheckThreadID(posts[0].Parent)
		if err != nil {
			return nil, fmt.Errorf("some error")
		}

		if pThread != threadID {
			return nil, fmt.Errorf("some error")
		}
	}

	query := `INSERT INTO posts(author, created, forum, message, parent, thread) VALUES `
	var args []interface{} //nolint:prealloc
	created := strfmt.DateTime(time.Now())

	for i, post := range posts {
		posts[i].Forum = forum
		posts[i].Thread = threadID
		posts[i].Created = created.String()

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6,
		)

		args = append(args, post.Author, created, forum, post.Message, post.Parent, threadID)
	}
	query = query[:len(query)-1]
	query += ` RETURNING id`

	rows, err := postRep.Conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error while inserting posts: %w", err)
	}
	defer rows.Close()

	var idx int
	for rows.Next() {
		err = rows.Scan(&posts[idx].ID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning %w", err)
		}

		idx++
	}
	return posts, nil
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
		rows, err = postRep.Conn.Query(context.Background(), query)
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
	var query string
	if since == "" {
		if desc {
			query = fmt.Sprintf(`SELECT id, parent, author, message, is_edited, forum, thread, created from posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 ORDER BY id DESC LIMIT %d)
				ORDER BY path[1] DESC, path, id;`, threadID, limit)
		} else {
			query = fmt.Sprintf(`SELECT id, parent, author, message, is_edited, forum, thread, created from posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 ORDER BY id LIMIT %d)
				ORDER BY path, id;`, threadID, limit)
		}
	} else {
		if desc {
			query = fmt.Sprintf(`SELECT id, parent, author, message, is_edited, forum, thread, created from posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 AND path[1] <
				(SELECT path[1] FROM posts WHERE id = %s) ORDER BY id DESC LIMIT %d) ORDER BY path[1] DESC, path, id;`,
				threadID, since, limit)
		} else {
			query = fmt.Sprintf(`SELECT id, parent, author, message, is_edited, forum, thread, created from posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 AND path[1] >
				(SELECT path[1] FROM posts WHERE id = %s) ORDER BY id ASC LIMIT %d) ORDER BY path, id;`,
				threadID, since, limit)
		}
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
