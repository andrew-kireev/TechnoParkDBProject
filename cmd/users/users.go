package main

import (
	forumHTTP "TechnoParkDBProject/internal/app/forum/delivery/http"
	forumRepositoru "TechnoParkDBProject/internal/app/forum/repository"
	forumUsecase "TechnoParkDBProject/internal/app/forum/usecase"
	threadHTTP "TechnoParkDBProject/internal/app/thread/delivery/http"
	threadRepositoru "TechnoParkDBProject/internal/app/thread/repository"
	threadUsecase "TechnoParkDBProject/internal/app/thread/usecase"
	userHTTP "TechnoParkDBProject/internal/app/user/delivery/http"
	"TechnoParkDBProject/internal/app/user/repository"
	"TechnoParkDBProject/internal/app/user/usecase"
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	//config := configs.NewConfig()

	router := router.New()

	dbpool, err := pgxpool.Connect(context.Background(),
		"host=localhost port=5432 dbname=db_project sslmode=disable",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRep := repository.NewUserRepository(dbpool)
	forumRep := forumRepositoru.NewUserRepository(dbpool)
	thredRep := threadRepositoru.NewThreadRepository(dbpool)

	userUsecase := usecase.NewUserUsecase(userRep)
	forumUsec := forumUsecase.NewForumUsecase(forumRep)
	thredUsec := threadUsecase.NewThreadUsecase(thredRep)

	userHTTP.NewUserHandler(router, userUsecase)
	forumHTTP.NewForumHandler(router, forumUsec, userUsecase)
	threadHTTP.NewThreadHandler(router, thredUsec)

	err = fasthttp.ListenAndServe(":5000", router.Handler)
	fmt.Println(err)
}
