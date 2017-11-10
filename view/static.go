package view

import (
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/static"
)

func StaticFile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.Response.Header.Add("Cache-Control", "max-age=2628000, private")

	switch string(ctx.Path()) {
	case "/static/gorum.css":
		ctx.SetContentType("text/css")
		static.WriteStyleSheet(ctx)
		return
	default:
		ctx.NotFound()
		return
	}
}
