package http

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/forum/models"
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/user"
	"TechnoParkDBProject/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ForumHandler struct {
	router       *router.Router
	forumUsecase forum.Usecase
	userUsecase  user.Usecase
}

func NewForumHandler(router *router.Router, forumUsecase forum.Usecase, userUsecase user.Usecase) *ForumHandler {
	forumHandler := &ForumHandler{
		router:       router,
		forumUsecase: forumUsecase,
		userUsecase:  userUsecase,
	}

	forumHandler.router.POST("/api/forum/create",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(forumHandler.CreateForumHandler)))
	forumHandler.router.GET("/api/forum/{forum_slug}/details",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(forumHandler.GetForumHandler)))
	forumHandler.router.GET("/api/forum/{forum_slug}/users",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(forumHandler.GetUsersByForumHandler)))

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
	user, err := handler.userUsecase.GetUserByNickname(newForum.UserNickname)
	if err != nil {
		resp := &responses.Response{
			Message: "Can't find user with nickname" + newForum.UserNickname,
		}
		body, err := resp.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}
	newForum.UserNickname = user.Nickname

	err = handler.forumUsecase.CreateForum(newForum)
	if err != nil {
		alredyExictedForum, err := handler.forumUsecase.GetForumBySlug(newForum.Slug)
		if err != nil {
			resp := &responses.Response{
				Message: "Can't find user with nickname" + newForum.UserNickname,
			}
			body, err := resp.MarshalJSON()
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}
			ctx.SetBody(body)
			ctx.SetStatusCode(http.StatusNotFound)
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

func (handler *ForumHandler) GetForumHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("forum_slug").(string)
	forum, err := handler.forumUsecase.GetForumBySlug(slug)
	if err != nil {
		resp := &responses.Response{
			Message: "Can't find forum with slug " + slug,
		}
		body, err := resp.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}
	body, err := forum.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(body)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *ForumHandler) GetUsersByForumHandler(ctx *fasthttp.RequestCtx) {
	forumSlug := ctx.UserValue("forum_slug").(string)
	desc := ctx.QueryArgs().GetBool("desc")
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 100
	}
	since := string(ctx.QueryArgs().Peek("since"))

	users, err := handler.forumUsecase.GetUsersByForum(forumSlug, since, limit, desc)
	if err != nil {
		resp := &responses.Response{
			Message: "Can't find users with forum slug " + forumSlug,
		}
		body, err := resp.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}
	body, err := json.Marshal(users)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(body)
	ctx.SetStatusCode(http.StatusOK)
}
