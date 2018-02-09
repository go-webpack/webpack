package main

import (
	"flag"
	"log"

	"github.com/go-webpack/webpack"
	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func homeIndex(ctx *iris.Context) {
	ctx.MustRender("home.html", struct{}{})
}

func main() {
	isDev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	webpack.Init(*isDev)
	view := view.HTML("./templates", ".html")
	view = view.Layout("layout.html")
	view = view.Funcs(map[string]interface{}{"asset": webpack.AssetHelper})

	app := iris.New()
	app.Adapt(view.Reload(*isDev))
	app.Adapt(httprouter.New())
	app.Get("/", homeIndex)

	log.Println("Iris demo app listening on http://localhost:3200")
	app.Listen(":3200")
}
