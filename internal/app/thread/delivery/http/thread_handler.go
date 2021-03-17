package http

import (
	"TechnoParkDBProject/internal/app/thread"
	"TechnoParkDBProject/internal/app/thread/models"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ThreadHandler struct {
	router        *router.Router
	threadUsecase thread.Usecase
}

func NewThreadHandler(router *router.Router, threadUsecase thread.Usecase) *ThreadHandler {
	threadHandler := &ThreadHandler{
		router:        router,
		threadUsecase: threadUsecase,
	}
	return threadHandler
}

func (handlre *ThreadHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	req := &models.Request{}
	err := req.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	thread := &models.Thread{
		Title:   "",
		Author:  "",
		Forum:   slug,
		Message: req.Message,
		Slug:    slug,
	}
	err = handlre.threadUsecase.CreateThread()
}
