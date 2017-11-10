package view

import (
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/model"
	"github.com/tiehuis/gorum/template"
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")

	bs, _ := model.GetAllBoards()

	page := &template.IndexPage{bs}
	template.WritePageTemplate(ctx, page)
}
