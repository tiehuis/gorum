//go:generate qtc -dir=template
//go:generate qtc -dir=static

package main

import (
	"flag"

	"github.com/buaazp/fasthttprouter"
	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"

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
	dbName := flag.String("sqlite-db", "file::memory:?mode=memory&cache=shared", "Database name")
	listenAddr := flag.String("listen-address", ":8080", "Server address")
	flag.Parse()

	model.Init(*dbName)
	model.Migrate()

	rlog.Info("Started server on", *listenAddr)

	view.InitMemCache()
	routes := routes().Handler
	routes = view.TimerHandler(routes)

	rlog.Critical(fasthttp.ListenAndServe(*listenAddr, routes))
}
