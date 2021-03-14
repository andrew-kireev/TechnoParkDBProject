package main

import (
	"TechnoParkDBProject/internal/app/user/models"
	"TechnoParkDBProject/internal/app/user/repository"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(),
		"host=localhost port=5432 dbname=db_project sslmode=disable",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRep := repository.NewUserRepository(dbpool)

	user := &models.User{
		Nickname: "some nick",
		FullName: "andrew kireev",
		About:    "oaoaoaoa",
		Email:    "a-kireev1989@mail.ru",
	}

	err = userRep.CreateUser(user)
	fmt.Println(err)
}
