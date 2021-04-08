package http

import (
	"TechnoParkDBProject/internal/app/middlware"
	"TechnoParkDBProject/internal/app/vote"
	"TechnoParkDBProject/internal/app/vote/models"
	"TechnoParkDBProject/internal/pkg/responses"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type VoteHandler struct {
	router      *router.Router
	voteUsecase vote.Usecase
}

func NewVoteHandler(router *router.Router, voteUsecase vote.Usecase) *VoteHandler {
	voteHandler := &VoteHandler{
		router:      router,
		voteUsecase: voteUsecase,
	}

	voteHandler.router.POST("/api/thread/{slug_or_id}/vote",
		middlware.ContentTypeJson(voteHandler.CreateVote))

	return voteHandler
}

func (handler *VoteHandler) CreateVote(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	vote := &models.Vote{}
	err := vote.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	thread, err := handler.voteUsecase.CreateNewVote(vote, slugOrID)
	if err != nil {
		response := responses.Response{Message: "Can't find thread with slug " + slugOrID}
		body, _ := response.MarshalJSON()
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}
	body, err := thread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
