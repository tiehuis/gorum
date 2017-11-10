package view

import (
	"time"

	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"
)

func TimerHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		h(ctx)
		end := time.Now()

		rlog.Debugf("%s '%s' took %v", ctx.Method(), ctx.Path(), end.Sub(start))
	}
}
