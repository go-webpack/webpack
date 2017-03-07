package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/qor/utils"
	"github.com/qor/render"
	"github.com/qor/wildcard_router"
	webpack "gopkg.in/webpack.v0"
)

func HomeIndex(ctx *gin.Context) {
	Render.Funcs(ViewHelpers()).Execute(
		"home_index",
		gin.H{},
		ctx.Request,
		ctx.Writer,
	)
}

var Render *render.Render

func init() {
	Render = render.New()
}

func ViewHelpers() map[string]interface{} {
	return map[string]interface{}{"asset": webpack.AssetHelper}
}

func main() {
	is_dev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	webpack.Init(*is_dev)

	mux := http.NewServeMux()

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.GET("/", HomeIndex)

	for _, path := range []string{"webpack"} {
		mux.Handle(fmt.Sprintf("/%s/", path), utils.FileServer(http.Dir("public")))
	}

	WildcardRouter := wildcard_router.New()
	WildcardRouter.MountTo("/", mux)
	WildcardRouter.AddHandler(router)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}
