package http

import (
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/user"
	"TechnoParkDBProject/internal/app/user/models"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserHandler struct {
	router      *fasthttprouter.Router
	userUsecase user.Useacse
}

func NewUserHandler(router *fasthttprouter.Router, userUsecase user.Useacse) *UserHandler {
	userHandler := &UserHandler{
		router:      router,
		userUsecase: userUsecase,
	}

	userHandler.router.POST("/user/:nickname/create",
		middlware.ContentTypeJson(userHandler.CreateUserHandler))
	userHandler.router.GET("/user/:nickname/profile",
		middlware.ContentTypeJson(userHandler.GetUserHandler))


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
		newUser, err := handler.userUsecase.GetUserByEmailOrNickname(newUser.Nickname, newUser.Email)
		if err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		body, err := newUser.MarshalJSON()
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
}

func (handler *UserHandler) GetUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := handler.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		body := fmt.Sprintf("{\n\"message\": \"Can't find user with nickname %v\n}", nickname)
		ctx.SetBody([]byte(body))
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
