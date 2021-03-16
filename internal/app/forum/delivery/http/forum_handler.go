package http

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/middlware"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ForumHandler struct {
	router       *fasthttprouter.Router
	forumUsecase forum.Usecase
}

func NewForumHandler(router *fasthttprouter.Router, forumUsecase forum.Usecase) *ForumHandler {
	forumHandler := &ForumHandler{
		router:       router,
		forumUsecase: forumUsecase,
	}

	forumHandler.router.POST("/api/forum/create",
		middlware.ContentTypeJson(forumHandler.CreateForumHandler))

	return forumHandler
}

func (handler *ForumHandler) CreateForumHandler(ctx *fasthttp.RequestCtx) {
	newForum := &models.Forum{}
	err := json.Unmarshal(ctx.PostBody(), newForum)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	err = handler.forumUsecase.CreateForum(newForum)
	if err != nil {
		alredyExictedForum, err := handler.forumUsecase.GetForumBySlug(newForum.Slug)
		if err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			body := fmt.Sprintf("{\n\"message\": \"Can't find user with nickname %v\n}",
				newForum.UserNickname)
			ctx.SetBody([]byte(body))
			return
		}
		fmt.Println(err)
		body, err := alredyExictedForum.MarshalJSON()
		if err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}
	body, err := newForum.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}
