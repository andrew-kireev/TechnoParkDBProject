package main

import (
	"TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/forum/repository"
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

	forumRep := repository.NewUserRepository(dbpool)
	forum := &models.Forum{
		Tittle:       "title",
		UserNickname: "andrew",
		Slug:         "http://some_url.com",
		Posts:        101,
		Threads:      5,
	}

	err = forumRep.CreateForum(forum)
	fmt.Println(err)
}
