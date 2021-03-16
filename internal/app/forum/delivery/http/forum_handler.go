package http

import (
	"TechnoParkDBProject/internal/app/user"
	"github.com/buaazp/fasthttprouter"
)

type ForumHandler struct {
	router      *fasthttprouter.Router
	forumUsecase user.Usecase
}

func NewForumHandler(router *fasthttprouter.Router, userUsecase user.Usecase) *ForumHandler {
	forumHandler := &ForumHandler{
		router:      router,
		forumUsecase: userUsecase,
	}

	return forumHandler
}
