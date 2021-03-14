package repository

import (
	"TechnoParkDBProject/internal/app/user/models"
	"context"
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