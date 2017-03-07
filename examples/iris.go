package main

import (
	"flag"
	"log"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	webpack "gopkg.in/webpack.v0"
)

func HomeIndex(ctx *iris.Context) {
	ctx.MustRender("home.html", struct{}{})
}

func main() {
	is_dev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	webpack.Init(*is_dev)
	view := view.HTML("./templates", ".html")
	view = view.Layout("layout.html")
	view = view.Funcs(map[string]interface{}{"asset": webpack.AssetHelper})

	app := iris.New()
	app.Adapt(view.Reload(*is_dev))
	app.Adapt(httprouter.New())
	app.Get("/", HomeIndex)

	log.Println("Iris demo app listening on http://localhost:3200")
	app.Listen(":3200")
}
