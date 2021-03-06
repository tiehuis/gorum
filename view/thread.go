package view

import (
	"fmt"
	"strconv"

	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/model"
	"github.com/tiehuis/gorum/template"
)

func Thread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")

	tpath := ctx.UserValue("thread").(string)
	tId, err := strconv.ParseInt(tpath, 10, 64)
	if err != nil {
		rlog.Debug("failed to parse thread", tpath, err)
		NotFound(ctx)
		return
	}

	t, err := model.GetPostById(tId)
	if err != nil {
		rlog.Debug("failed to get post", tId, err)
		NotFound(ctx)
		return
	}

	b, err := model.GetBoardById(t.BoardParentId)
	if err != nil {
		rlog.Debug("failed to get board with id", t.BoardParentId, err)
		NotFound(ctx)
		return
	}

	if t.ThreadParentId.Valid {
		rp := fmt.Sprintf("/board/%s/%v#%v", b.Code, t.ThreadParentId.Int64, t.Id)
		rlog.Debug("Performing Thread redirect to", rp)
		ctx.Redirect(rp, fasthttp.StatusFound)
		return
	}

	board := ctx.UserValue("board").(string)
	if b.Code != board {
		rlog.Debugf("requested post with wrong board; found %s but post is on %s", board, b.Code)
		NotFound(ctx)
		return
	}

	ts, err := t.GetParentThread()
	if err != nil {
		rlog.Debug("failed to get parent thread id of post", t.Id, err)
		NotFound(ctx)
		return
	}

	page := &template.ThreadPage{b, ts}
	template.WritePageTemplate(ctx, page)
}

func ThreadPost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	a := ctx.PostArgs()

	if isRateLimited(ctx) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "you have just posted something, try again soon")
		return
	}

	board := ctx.UserValue("board").(string)
	b, err := model.GetBoardByCode(board)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "could not find board")
		return
	}

	tis := ctx.UserValue("thread").(string)
	tid, err := strconv.ParseInt(tis, 10, 64)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "could not parse thread id")
		return
	}

	t, err := model.GetPostById(tid)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "could not find post id")
		return
	}

	if t.ThreadParentId.Valid {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "post id isn't a top level thread")
		return
	}

	if t.ArchivedAt.Valid {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "cannot post to an archived thread")
		return
	}

	content := string(a.Peek("content"))
	if len(content) == 0 || len(content) >= 2000 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "content must be between 1 and 2000 bytes inclusive")
		return
	}

	pc, err := t.GetParentThreadCount()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "could not determine thread count")
		return
	}

	// If two concurrent requests both pass this at a similar time we can have
	// more than 200 posts in a thread. This isn't a big deal so we don't worry
	// about it.
	if pc >= 200 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		GeneralError(ctx, "thread post limit has been reached")
		return
	}

	rlog.Debugf("Creating new post: %d %d", tid, len(content))
	nid, err := model.CreatePost(model.PostW{tid, content})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		GeneralError(ctx, "failed to create post right now")
		rlog.Debug("Failed to create post:", err)
		return
	}

	startRateLimit(ctx)

	rp := fmt.Sprintf("/board/%s/%d/#%d", b.Code, tid, nid)
	rlog.Debug("Performing ThreadPost redirect to", rp)
	ctx.Redirect(rp, fasthttp.StatusFound)
}
