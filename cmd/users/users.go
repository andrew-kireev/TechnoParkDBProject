package main

import (
	userHTTP "TechnoParkDBProject/internal/app/user/delivery/http"
	"TechnoParkDBProject/internal/app/user/repository"
	"TechnoParkDBProject/internal/app/user/usecase"
	"context"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	//config := configs.NewConfig()

	router := fasthttprouter.New()

	dbpool, err := pgxpool.Connect(context.Background(),
		"host=localhost port=5432 dbname=db_project sslmode=disable",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRep := repository.NewUserRepository(dbpool)
	userUsecase := usecase.NewUserUsecase(userRep)

	userHTTP.NewUserHandler(router, userUsecase)

	//serv := server.NewServer(router, userHandler)

	err = fasthttp.ListenAndServe(":7777", router.Handler)
	fmt.Println(err)
}
