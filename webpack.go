package webpack

import (
	"errors"
	"html/template"
	"log"
	"strings"

	"github.com/go-webpack/webpack/helper"
	"github.com/go-webpack/webpack/reader"
)

var DevHost = "localhost:3808"
var FsPath = "public/webpack"
var WebPath = "webpack"
var Plugin = "stats"
var IgnoreMissing = true
var Verbose = true

var isDev = false
var initDone = false
var preloadedAssets map[string][]string

func readManifest() (map[string][]string, error) {
	return reader.Read(Plugin, DevHost, FsPath, WebPath, isDev)
}

// Init Set current environment and preload manifest
func Init(dev bool) {
	var err error
	isDev = dev
	if isDev {
		// Try to preload manifest, so we can show an error if webpack-dev-server is not running
		_, err = readManifest()
	} else {
		preloadedAssets, err = readManifest()
	}
	if err != nil {
		log.Println("go-webpack: Error loading asset manifest", err)
	}
	initDone = true
}

func AssetHelper(key string) (template.HTML, error) {
	var err error

	if !initDone {
		return "", errors.New("Please call webpack.Init() first (see readme)")
	}

	var assets map[string][]string
	if isDev {
		assets, err = readManifest()
		if err != nil {
			return template.HTML(""), err
		}
	} else {
		assets = preloadedAssets
	}

	parts := strings.Split(key, ".")
	kind := parts[len(parts)-1]
	//log.Println("showing assets:", key, parts, kind)

	v, ok := assets[key]
	if !ok {
		message := "go-webpack: Asset file '" + key + "' not found in manifest"
		if Verbose {
			log.Printf("%s. Manifest contens: %+v", message, assets)
		}
		if IgnoreMissing {
			return template.HTML(""), nil
		} else {
			return template.HTML(""), errors.New(message)
		}
	}

	buf := []string{}
	for _, s := range v {
		if strings.HasSuffix(s, "."+kind) {
			buf = append(buf, helper.AssetTag(kind, s))
		} else {
			log.Println("skip asset", s, ": bad type")
		}
	}
	return template.HTML(strings.Join(buf, "\n")), nil
}
