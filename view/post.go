package view

import (
	"fmt"
	"strconv"

	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/model"
)

func Post(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		rlog.Debug("failed to parse id", id, err)
		NotFound(ctx)
		return
	}

	p, err := model.GetPostById(postId)
	if err != nil {
		rlog.Debug("failed to find post", err)
		NotFound(ctx)
		return
	}

	b, err := model.GetBoardById(p.BoardParentId)
	if err != nil {
		rlog.Debug("failed to find parent board", b)
		NotFound(ctx)
		return
	}

	var rp string
	if p.ThreadParentId == 0 {
		rp = fmt.Sprintf("/board/%s/%v", b.Code, p.Id)
	} else {
		rp = fmt.Sprintf("/board/%s/%v#%v", b.Code, p.ThreadParentId, p.Id)
	}

	rlog.Debug("Performing Post redirect to", rp)
	ctx.Redirect(rp, fasthttp.StatusFound)
}
