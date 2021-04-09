package repository

import (
	"TechnoParkDBProject/internal/app/user/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(con *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Conn: con,
	}
}

func (userRep *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (nickname, fullname, about, email)
			VALUES ($1, $2, $3, $4)`

	_, err := userRep.Conn.Exec(context.Background(), query, user.Nickname, user.FullName,
		user.About, user.Email)

	return err
}

func (userRep *UserRepository) GetUserByEmailOrNickname(nickname, email string) ([]*models.User, error) {
	query := `SELECT nickname, fullname, about, email from users
			where nickname = $1 or email = $2`
	//user := &models.User{}
	users := make([]*models.User, 0)

	rows, err := userRep.Conn.Query(context.Background(), query,
		nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		_ = rows.Scan(&user.Nickname, &user.FullName, &user.About,
			&user.Email)

		users = append(users, user)
	}
	return users, nil
}

func (userRep *UserRepository) GetUserByNickname(nickname string) (*models.User, error) {
	query := `SELECT nickname, fullname, about, email from users
			where nickname = $1`
	user := &models.User{}

	err := userRep.Conn.QueryRow(context.Background(), query,
		nickname).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRep *UserRepository) UpdateUserInformation(user *models.User) error {
	query := `UPDATE users
SET fullname = (CASE
                    WHEN LTRIM($1) = '' THEN fullname
                    ELSE $1 END)
  , about=(CASE
               WHEN $2 = '' THEN about
               ELSE $2 END)
  , email= (CASE
                    WHEN LTRIM($3) = '' THEN email
                    ELSE LTRIM($3)
    END)
where nickname = $4`

	_, err := userRep.Conn.Exec(context.Background(), query, user.FullName,
		user.About, user.Email, user.Nickname)

	if err != nil {
		return err
	}
	return nil
}

func (userRep *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT nickname, fullname, about, email from users
			where email = $1`
	user := &models.User{}

	err := userRep.Conn.QueryRow(context.Background(), query,
		email).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRep *UserRepository) DeleteAll() error {
	query := `TRUNCATE TABLE users_to_forums CASCADE;
			TRUNCATE TABLE votes CASCADE;
			TRUNCATE TABLE posts CASCADE;
			TRUNCATE TABLE threads CASCADE;
			TRUNCATE TABLE forum CASCADE;
			TRUNCATE TABLE users CASCADE;`

	_, err := userRep.Conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (userRep *UserRepository) GetStatus() (*models.Status, error){
	queryUser := `SELECT COUNT(*) as user_count FROM users;`
	queryForum := `SELECT COUNT(*) as forum_count FROM forum;`
	queryThread := `SELECT COUNT(*) as thread_count FROM threads;`
	queryPost := `SELECT COUNT(*) as post_count FROM posts;`

	status := &models.Status{}

	err := userRep.Conn.QueryRow(context.Background(), queryUser).Scan(&status.User)
	if err != nil {
		return nil, err
	}
	err = userRep.Conn.QueryRow(context.Background(), queryForum).Scan(&status.Forum)
	if err != nil {
		return nil, err
	}
	err = userRep.Conn.QueryRow(context.Background(), queryThread).Scan(&status.Thread)
	if err != nil {
		return nil, err
	}
	err = userRep.Conn.QueryRow(context.Background(), queryPost).Scan(&status.Post)
	if err != nil {
		return nil, err
	}

	return status, nil
}
