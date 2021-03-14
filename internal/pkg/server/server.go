package server

import (
	userHTTP "TechnoParkDBProject/internal/app/user/delivery/http"
	fasthttprouter "github.com/buaazp/fasthttprouter"
)

type Server struct {
	Usecase *userHTTP.UserHandler
	Router *fasthttprouter.Router
}

func NewServer(router *fasthttprouter.Router, userHandler *userHTTP.UserHandler) *Server {
	serv := &Server{
		Usecase: userHandler,
		Router: router,
	}

	return serv
}
