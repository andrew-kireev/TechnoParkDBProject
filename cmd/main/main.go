package main

import (
	forumHTTP "TechnoParkDBProject/internal/app/forum/delivery/http"
	forumRepository "TechnoParkDBProject/internal/app/forum/repository"
	forumUsecase "TechnoParkDBProject/internal/app/forum/usecase"
	postHTTP "TechnoParkDBProject/internal/app/posts/delivery/http"
	postRepository "TechnoParkDBProject/internal/app/posts/repository"
	postUsecase "TechnoParkDBProject/internal/app/posts/usecase"
	threadHTTP "TechnoParkDBProject/internal/app/thread/delivery/http"
	threadRepositoru "TechnoParkDBProject/internal/app/thread/repository"
	threadUsecase "TechnoParkDBProject/internal/app/thread/usecase"
	userHTTP "TechnoParkDBProject/internal/app/user/delivery/http"
	"TechnoParkDBProject/internal/app/user/repository"
	"TechnoParkDBProject/internal/app/user/usecase"
	voteHTTP "TechnoParkDBProject/internal/app/vote/delivery/http"
	voteRepository "TechnoParkDBProject/internal/app/vote/repository"
	voteUsecase "TechnoParkDBProject/internal/app/vote/usecase"
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

	// "host=localhost port=5432 dbname=db_project sslmode=disable"
	// host=localhost port=5432 user=andrewkireev dbname=db_forum password=password sslmode=disable
	dbpool, err := pgxpool.Connect(context.Background(),
		"host=localhost port=5432 dbname=db_project sslmode=disable",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRep := repository.NewUserRepository(dbpool)
	forumRep := forumRepository.NewUserRepository(dbpool)
	thredRep := threadRepositoru.NewThreadRepository(dbpool)
	postsRep := postRepository.NewPostsRepository(dbpool)
	voteRep := voteRepository.NewVoteRepository(dbpool)

	userUsecase := usecase.NewUserUsecase(userRep)
	forumUsec := forumUsecase.NewForumUsecase(forumRep)
	thredUsec := threadUsecase.NewThreadUsecase(thredRep)
	postUse := postUsecase.NewPostsUsecase(postsRep, thredRep, forumRep, userRep)
	voteUse := voteUsecase.NewVoteUsecase(voteRep, thredRep)

	userHTTP.NewUserHandler(router, userUsecase)
	forumHTTP.NewForumHandler(router, forumUsec, userUsecase)
	threadHTTP.NewThreadHandler(router, thredUsec, forumUsec)
	postHTTP.NewPostsHandler(router, postUse)
	voteHTTP.NewVoteHandler(router, voteUse)

	err = fasthttp.ListenAndServe(":5000", router.Handler)
	fmt.Println(err)
}
