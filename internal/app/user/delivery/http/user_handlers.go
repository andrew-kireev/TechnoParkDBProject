package http

import (
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/user"
	"TechnoParkDBProject/internal/app/user/models"
	"TechnoParkDBProject/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserHandler struct {
	router      *fasthttprouter.Router
	userUsecase user.Usecase
}

func NewUserHandler(router *fasthttprouter.Router, userUsecase user.Usecase) *UserHandler {
	userHandler := &UserHandler{
		router:      router,
		userUsecase: userUsecase,
	}

	userHandler.router.POST("/api/user/:nickname/create",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(userHandler.CreateUserHandler)))
	userHandler.router.GET("/api/user/:nickname/profile",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(userHandler.GetUserHandler)))
	userHandler.router.POST("/api/user/:nickname/profile",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(userHandler.UpdateUserHandler)))
	userHandler.router.POST("/api/service/clear",
		middlware.LoggingMiddleware(middlware.ContentTypeJson(userHandler.DeleteAllHandler)))
	return userHandler
}

func (handler *UserHandler) CreateUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	newUser := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), newUser)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Nickname = nickname
	err = handler.userUsecase.CreateUser(newUser)
	if err != nil {
		newUsers, err := handler.userUsecase.GetUserByEmailOrNickname(newUser.Nickname, newUser.Email)
		if err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(newUsers)
		if err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
		}
		ctx.SetBody(body)
		ctx.SetStatusCode(http.StatusConflict)
		return
	}

	body, err := newUser.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	ctx.SetBody(body)
	ctx.SetStatusCode(http.StatusCreated)
}

func (handler *UserHandler) GetUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := handler.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		resp := &responses.Response{
			Message: "Can't find user with nickname" + nickname,
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

	body, err := user.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	ctx.SetBody(body)
}

func (handler *UserHandler) UpdateUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	newUser := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), newUser)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Nickname = nickname
	fmt.Println(newUser)

	us, err := handler.userUsecase.GetUserByEmail(newUser.Email)
	if err == nil && us.Nickname != newUser.Nickname {
		resp := &responses.Response{
			Message: "Can't find user with nickname" + nickname,
		}
		body, err := resp.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetBody(body)
		ctx.SetStatusCode(http.StatusConflict)
		return
	}
	oldUser, err := handler.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		resp := &responses.Response{
			Message: "Can't find user with nickname" + nickname,
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

	err = handler.userUsecase.UpdateUserInformation(newUser)

	if newUser.FullName == "" {
		newUser.FullName = oldUser.FullName
	}
	if newUser.About == "" {
		newUser.About = oldUser.About
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}

	body, err := newUser.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
	}
	ctx.SetBody(body)
}

func (handler *UserHandler) DeleteAllHandler(ctx *fasthttp.RequestCtx) {
	err := handler.userUsecase.DeleteAll()
	if err != nil {
		fmt.Println(err)
	}

	ctx.SetStatusCode(http.StatusOK)
}
