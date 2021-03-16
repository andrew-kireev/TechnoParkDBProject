package repository

import "github.com/jackc/pgx/v4/pgxpool"

type ForumRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(con *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		Conn: con,
	}
}
