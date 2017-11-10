package view

import (
	"fmt"

	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/model"
	"github.com/tiehuis/gorum/template"
)

func Board(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")

	path := ctx.UserValue("board").(string)
	b, err := model.GetBoardByCode(path)
	if err != nil {
		rlog.Debug("failed to find board", path, err)
		ctx.NotFound()
		return
	}

	ps, err := b.GetAllPosts()
	if err != nil {
		rlog.Debug("failed to get posts for board", path, err)
		ctx.NotFound()
		return
	}

	page := &template.BoardPage{Board: b, Threads: ps}
	template.WritePageTemplate(ctx, page)
}

func BoardPost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	a := ctx.PostArgs()

	board := ctx.UserValue("board").(string)
	b, err := model.GetBoardByCode(board)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("could not find board"))
		return
	}

	content := string(a.Peek("content"))
	if len(content) == 0 || len(content) >= 2000 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("content must be between 1 and 2000 bytes inclusive"))
		return
	}

	rlog.Debugf("Creating new thread: %d %d", b.Id, len(content))
	_, err = model.CreateThread(model.ThreadW{b.Id, content})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("failed to create post right now"))
		return
	}

	rp := fmt.Sprintf("/board/%s", b.Code)
	rlog.Debug("Performing BoardPost redirect to", rp)
	ctx.Redirect(rp, fasthttp.StatusFound)
}
