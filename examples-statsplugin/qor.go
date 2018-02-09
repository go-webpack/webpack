package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webpack/webpack"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/qor/utils"
	"github.com/qor/render"
	"github.com/qor/wildcard_router"
)

var renderer *render.Render

func init() {
	renderer = render.New(&render.Config{
		//ViewPaths:     []string{"app/views"},
		//DefaultLayout: "application", // default value is application
		FuncMapMaker: func(*render.Render, *http.Request, http.ResponseWriter) template.FuncMap {
			return viewHelpers()
		},
	})
}

func homeIndex(ctx *gin.Context) {
	// Alternative (without FuncMapMaker):
	//renderer.Funcs(viewHelpers()).Execute(

	renderer.Execute(
		"home_index",
		gin.H{},
		ctx.Request,
		ctx.Writer,
	)
}

func viewHelpers() map[string]interface{} {
	return map[string]interface{}{"asset": webpack.AssetHelper}
}

func main() {
	isDev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	webpack.Init(*isDev)

	mux := http.NewServeMux()

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.GET("/", homeIndex)

	for _, path := range []string{"webpack"} {
		mux.Handle(fmt.Sprintf("/%s/", path), utils.FileServer(http.Dir("public")))
	}

	WildcardRouter := wildcard_router.New()
	WildcardRouter.MountTo("/", mux)
	WildcardRouter.AddHandler(router)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}
