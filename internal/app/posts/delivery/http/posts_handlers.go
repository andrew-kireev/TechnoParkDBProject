package http

import (
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/posts"
	"TechnoParkDBProject/internal/app/posts/models"
	"TechnoParkDBProject/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
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

	postsHandler.router.POST("/api/thread/{slug}/create",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(postsHandler.CreatePost)))
	postsHandler.router.GET("/api/thread/{slug_or_id}/posts", middlware.ContentTypeJson(postsHandler.GetPosts))
	postsHandler.router.GET("/api/post/{id}/details", middlware.ContentTypeJson(postsHandler.GetPostHandler))
	postsHandler.router.POST("/api/post/{id}/details", middlware.ContentTypeJson(postsHandler.UpdatePost))

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
		switch {
		case err.Error() == "no user":
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusNotFound)
			resp := responses.Response{Message: "Can't find user"}
			body, _ := resp.MarshalJSON()
			ctx.SetBody(body)
			return
		case err.Error() == "no thread":
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusNotFound)
			resp := responses.Response{Message: "Can't find thread or forum"}
			body, _ := resp.MarshalJSON()
			ctx.SetBody(body)
			return
		case err.Error() == "no posts":
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusCreated)
			emptyPosts := make([]*models.Post, 0)
			body, _ := json.Marshal(emptyPosts)
			ctx.SetBody(body)
			return
		default:
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusConflict)
			resp := responses.Response{Message: "Parent post was created in another thread"}
			body, _ := resp.MarshalJSON()
			ctx.SetBody(body)
			return
		}
	}

	fmt.Println("fsdfds")
	ctx.SetStatusCode(http.StatusCreated)
	body, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(body)
}

func (postHandler *PostsHandler) GetPosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	since := string(ctx.QueryArgs().Peek("since"))
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		fmt.Println(err)
	}
	sort := string(ctx.QueryArgs().Peek("sort"))
	desc := ctx.QueryArgs().GetBool("desc")
	posts, err := postHandler.postUsecase.GetPosts(sort, since, slugOrID, limit, desc)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusNotFound)
		resp := responses.Response{Message: "Can't threads with forum slug " + slugOrID}
		body, _ := resp.MarshalJSON()
		ctx.SetBody(body)
		return
	}
	body, err := json.Marshal(posts)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (handler *PostsHandler) GetPostHandler(ctx *fasthttp.RequestCtx) {
	postID, err := strconv.Atoi(ctx.UserValue("id").(string))

	relatedArr := string(ctx.QueryArgs().Peek("related"))
	relatedStrs := strings.Split(relatedArr, ",")
	for len(relatedStrs) < 3 {
		relatedStrs = append(relatedStrs, "")
	}
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	postResponse, err := handler.postUsecase.GetPost(postID, relatedStrs)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusNotFound)
		resp := responses.Response{Message: fmt.Sprintf("Can't find post with id %d", postID)}
		body, _ := resp.MarshalJSON()
		ctx.SetBody(body)
		return
	}

	body, err := postResponse.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (handler *PostsHandler) UpdatePost(ctx *fasthttp.RequestCtx) {
	postID, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	post := &models.Post{}
	err = post.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	post.ID = postID
	post, err = handler.postUsecase.UpdatePost(post)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusNotFound)
		resp := responses.Response{Message: fmt.Sprintf("Can't find post with id: %d", postID)}
		body, _ := resp.MarshalJSON()
		ctx.SetBody(body)
		return
	}
	body, err := post.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
