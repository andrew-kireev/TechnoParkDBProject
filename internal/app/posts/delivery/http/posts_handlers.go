package http

import (
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/posts"
	"TechnoParkDBProject/internal/app/posts/models"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
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

	postsHandler.router.POST("/api/thread/{slug}/create", middlware.ContentTypeJson(postsHandler.CreatePost))

	return postsHandler
}

func (postHandler *PostsHandler) CreatePost(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	posts := make([]*models.Post, 0)

	err := json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	posts, err = postHandler.postUsecase.CreatePost(posts, slug)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if posts == nil {
		return
	}
	for _, post := range posts {
		fmt.Println(post)
	}
	ctx.SetStatusCode(http.StatusCreated)
	body, err := json.Marshal(posts)
	ctx.SetBody(body)
}
