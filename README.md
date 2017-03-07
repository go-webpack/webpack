## Introduction

This module allows proper integration with webpack, with support for proper assets reloading in development and asset hashes for production caching.

This module is compatible with both webpack 2.0 and 1.0. Example config file is for 2.0.

#### Usage with QOR
##### main.go
```golang
import (
  ...
	webpack "gopkg.in/webpack.v0"
)
func main() {
	is_dev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	webpack.Init(*is_dev)
  ...
}
```

##### controller.go
```golang
package controllers

import (
	"github.com/qor/render"
	"github.com/gin-gonic/gin"
	webpack "gopkg.in/webpack.v0"
)

var Render *render.Render

func init() {
	Render = render.New()
}

func ViewHelpers() map[string]interface{} {
	return map[string]interface{}{"asset": webpack.AssetHelper}
}

func HomeIndex(ctx *gin.Context) {
	Render.Funcs(ViewHelpers()).Execute(
		"home_index",
		gin.H{},
		ctx.Request,
		ctx.Writer,
	)
}
```

##### layouts/application.tmpl

```html
<!doctype html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    {{ asset "vendor.css" }}
    {{ asset "application.css" }}
  </head>
  <body>
    <div class="page-wrap">
      {{render .Template}}
    </div>
    {{ asset "vendor.js" }}
    {{ asset "application.js" }}
  </body>
</html>
```

#### Usage with Iris

##### main.go

```golang
import (
    webpack "gopkg.in/webpack.v0"
    iris "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func main() {
    is_dev := flag.Bool("dev", false, "development mode")
    flag.Parse()
    webpack.Init(*is_dev)
    view := view.HTML("./templates", ".html")
    view = view.Layout("layout.html")
    view = view.Funcs(map[string]interface{}{"asset": webpack.AssetHelper})
    app.Adapt(view.Reload(*is_dev))

    app.Adapt(httprouter.New())
}
```

##### templates/layout.html
```html
<!DOCTYPE HTML>
<html lang="en" >
<head>
<meta charset="UTF-8">
<title></title>
{{ asset "vendor.css" }}
{{ asset "application.css" }}
</head>
<body>
{{ yield }}
{{ asset "vendor.js" }}
{{ asset "application.js" }}
```

#### Usage with other frameworks

- Configure webpack to serve manifest.json via StatsPlugin
- Call ```webpack.Init()``` to set development or production mode.
- Add webpack.AssetHelper to your template functions.
- Call helper function with the name of your asset

Use webpack.config.js (and package.json) from this repo or create your own.

The only thing that must be present in your webpack config is StatsPlugin which is required to serve assets the proper way with hashes, etc.

Your assets is expected to be at public/webpack and your dev server at http://localhost:3808

When run with -dev flag, webpack asset manifest is loaded from http://localhost:3808/webpack/manifest.json, and updated automatically on every request. When running in production from public/webpack/manifest.json and is persistently cached in memory for performance reasons.

#### Running examples

```
cd examples
yarn install # or npm install
./node_modules/.bin/webpack-dev-server --config webpack.config.js --hot --inline
go get
go run iris.go -dev
go run qor.go -dev
```

#### Compiling assets for production

```
NODE_ENV=production ./node_modules/.bin/webpack --config webpack.config.js
```

#### License

Copyright (c) 2017 glebtv

MIT License


