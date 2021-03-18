package http

import (
	"TechnoParkDBProject/internal/app/posts"
	"github.com/fasthttp/router"
)

type PostsHandler struct {
	router      *router.Router
	postUsecase posts.Usecase
}

func NewPostsHandler(router *router.Router, usecase posts.Usecase) *PostsHandler {
	postsHandler := &PostsHandler{
		router:      router,
		postUsecase: usecase,
	}

	return postsHandler
}


//func (postHandler *PostsHandler) Create
