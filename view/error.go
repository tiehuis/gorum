package view

import (
	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/template"
)

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")

	page := &template.NotFoundPage{}
	template.WritePageTemplate(ctx, page)
}

func InternalError(ctx *fasthttp.RequestCtx, state interface{}) {
	ctx.SetContentType("text/html; charset=utf-8")
	rlog.Errorf("error handler:", state)

	page := &template.InternalErrorPage{state}
	template.WritePageTemplate(ctx, page)
}
