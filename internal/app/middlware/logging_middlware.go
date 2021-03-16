package middlware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func LoggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		msg := fmt.Sprintf("URL: %s, METHOD: %s, REMOTE_ADDR%s",
			ctx.URI(), ctx.Method(), ctx.RemoteAddr())
		logrus.Info(msg)
		next(ctx)
	}
}
