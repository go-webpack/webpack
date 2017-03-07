#### Introduction

This module allows proper integration with webpack, with support for asset hashes for production caching.

This module is compatible with both webpack 2.0 and 1.0. Example config file is for 2.0.

#### Usage with 

#### Usage with Iris

##### main.go
```
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
```
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


#### License

Copyright (c) 2017 glebtv

MIT License
