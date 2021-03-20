package http

import (
	"TechnoParkDBProject/internal/app/forum"
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/thread"
	"TechnoParkDBProject/internal/app/thread/models"
	"TechnoParkDBProject/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ThreadHandler struct {
	router        *router.Router
	threadUsecase thread.Usecase
	forumsUseacse forum.Usecase
}

func NewThreadHandler(router *router.Router, threadUsecase thread.Usecase,
	forumUsecase forum.Usecase) *ThreadHandler {
	threadHandler := &ThreadHandler{
		router:        router,
		threadUsecase: threadUsecase,
		forumsUseacse: forumUsecase,
	}

	router.POST("/api/forum/{slug}/create", middlware.ContentTypeJson(threadHandler.CreateThread))
	router.GET("/api/forum/{slug}/threads", middlware.ContentTypeJson(threadHandler.GetThreads))
	return threadHandler
}

func (handler *ThreadHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	thread := &models.Thread{Forum: slug}

	err := thread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	oldThread, err := handler.threadUsecase.FindThreadBySlug(thread.Slug)
	if err == nil && oldThread.Slug != "" {
		body, err := oldThread.MarshalJSON()
		if err != nil {
			fmt.Println(body)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}
	forum, err := handler.forumsUseacse.GetForumBySlug(thread.Forum)
	if err != nil {
		response := responses.Response{Message: "Can't find thread with slug " + slug}
		body, _ := response.MarshalJSON()
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}
	thread.Forum = forum.Slug
	thread, err = handler.threadUsecase.CreateThread(thread)
	if err != nil {
		fmt.Println(err)
		oldThread, errFind := handler.threadUsecase.FindThreadBySlug(slug)
		if errFind != nil {
			response := responses.Response{Message: "Can't find thread with slug " + slug}
			body, _ := response.MarshalJSON()
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(body)
			return
		}
		body, _ := oldThread.MarshalJSON()
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}
	body, err := thread.MarshalJSON()
	if err != nil {
		fmt.Println(body)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (handler *ThreadHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit, err := strconv.Atoi(string(ctx.FormValue("limit")))
	if err != nil {
		fmt.Println(err)
	}
	since := string(ctx.QueryArgs().Peek("since"))
	if err != nil {
		fmt.Println(err)
	}
	desc := string(ctx.FormValue("desc"))
	if err != nil {
		fmt.Println(err)
	}

	threads, err := handler.threadUsecase.GetThreadsByForumSlug(slug, since, desc, limit)
	if err != nil {
		fmt.Println("tut" + err.Error())
		ctx.SetStatusCode(http.StatusNotFound)
		resp := responses.Response{Message: "Can't threads with forum slug " + slug}
		body, _ := resp.MarshalJSON()
		ctx.SetBody(body)
		return
	}

	if len(threads) == 0 {
		_, err = handler.forumsUseacse.GetForumBySlug(slug)
		if err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			resp := responses.Response{Message: "Can't threads with forum slug " + slug}
			body, _ := resp.MarshalJSON()
			ctx.SetBody(body)
			return
		}
	}
	body, err := json.Marshal(threads)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
