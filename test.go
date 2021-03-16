package main

import (
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

	forumRep := repository.NewUserRepository(dbpool)

	err = forumRep.DeleteAll()
	fmt.Println(err)
}
