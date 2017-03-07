
#### Usage with Iris

##### main.go
```
import (
    "github.com/go-webpack/webpack"
    iris "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func main() {
    is_dev := flag.Bool("dev", false, "development mode")
    flag.Parse()
    wp.Init(*is_dev)
    view := view.HTML("./templates", ".html")
    view = view.Layout("layout.html")
    view = view.Funcs(map[string]interface{}{"asset": wp.AssetHelper})
    app.Adapt(view.Reload(*is_dev))

    app.Adapt(httprouter.New())
}
```

##### templates/layout.html
```
<!DOCTYPE HTML>
<html lang="en" >
<head>
<meta charset="UTF-8">
<title></title>
{{ asset "application.css" }}
</head>
<body>
{{ yield }}
{{ asset "application.js" }}
```

##### Usage with ... 

TODO

Use webpack.config.js (and package.json) from this repo or create your own.

The only thing that must be present in your webpack config is StatsPlugin which is required to serve assets the proper way with hashes, etc.

Your assets is expected to be at public/webpack and your dev server at http://localhost:3808

When run with -dev flag, webpack asset manifest is loaded from http://localhost:3808/webpack/manifest.json, and updated automatically on every request. When running in production from public/webpack/manifest.json and is persistently cached in memory for performance reasons.


#### License

Copyright (c) 2017 glebtv

MIT License
