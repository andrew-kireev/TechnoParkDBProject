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

func (userRep *UserRepository) GetUserByEmailOrNickname(nickname, email string) (*models.User, error) {
	query := `SELECT nickname, fullname, about, email from users
			where nickname = $1 or email = $2`
	user := &models.User{}

	err := userRep.Conn.QueryRow(context.Background(), query,
		nickname, email).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
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
	query := `UPDATE users SET fullname = $1, about = $2, email = $3
			where nickname = $4`

	res, err := userRep.Conn.Exec(context.Background(), query, user.FullName,
		user.About, user.Email, user.Nickname)
	fmt.Println(res)

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
