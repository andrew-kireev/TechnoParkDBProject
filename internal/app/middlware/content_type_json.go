package middlware

import (
	"github.com/valyala/fasthttp"
)

func ContentTypeJson(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		next(ctx)
	}
}
