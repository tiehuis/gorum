//go:generate qtc -dir=template
//go:generate qtc -dir=static

package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

	"github.com/tiehuis/gorum/config"
	"github.com/tiehuis/gorum/model"
	"github.com/tiehuis/gorum/view"
)

func routes() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/", view.Index)
	router.GET("/board/:board", view.Board)
	router.POST("/board/:board", view.BoardPost)
	router.GET("/board/:board/:thread", view.Thread)
	router.POST("/board/:board/:thread", view.ThreadPost)
	router.GET("/post/:id", view.Post)

	router.GET("/static/*path", view.StaticFile)

	router.NotFound = view.NotFound
	router.PanicHandler = view.InternalError
	return router
}

func main() {
	config.Init()

	model.Init()
	model.Migrate()
	model.PrepareQueries()

	rlog.Info("Started server on", config.HttpAddress)

	view.InitMemCache()
	routes := routes().Handler
	routes = view.TimerHandler(routes)

	rlog.Critical(fasthttp.ListenAndServe(config.HttpAddress, routes))
}
