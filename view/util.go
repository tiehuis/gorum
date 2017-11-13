package view

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"
)

var memcache *cache.Cache

func TimerHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		h(ctx)
		end := time.Now()

		rlog.Debugf("%s '%s' took %v", ctx.Method(), ctx.Path(), end.Sub(start))
	}
}

func InitMemCache() {
	memcache = cache.New(time.Minute, 2*time.Minute)
}

// This is a very rudimentary set of checks and should not be relied on.
//
// A anti-spam system such as reCaptcha would ideally be used in conjunction.
// I do not want to add a proper session management scheme server-side.
func isRateLimited(ctx *fasthttp.RequestCtx) bool {
	id := string(ctx.ConnID())
	_, found := memcache.Get(id)
	return found
}

func startRateLimit(ctx *fasthttp.RequestCtx) {
	id := string(ctx.ConnID())
	memcache.SetDefault(id, true)
}
