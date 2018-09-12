package webpack

import (
	"errors"
	"html/template"
	"log"
	"strings"

	"github.com/go-webpack/webpack/helper"
	"github.com/go-webpack/webpack/reader"
)

// DevHost webpack-dev-server host:port
var DevHost = "localhost:3808"

// FsPath filesystem path to public webpack dir
var FsPath = "./public/webpack"

// WebPath http path to public webpack dir
var WebPath = "webpack"

// Plugin webpack plugin to use, can be stats or manifest
var Plugin = "deprecated-stats"

// IgnoreMissing ignore assets missing on manifest or fail on them
var IgnoreMissing = true

// Verbose error messages to console (even if error is ignored)
var Verbose = true

var isDev = false
var initDone = false
var preloadedAssets map[string][]string

type Config struct {
	// DevHost webpack-dev-server host:port
	DevHost string
	// FsPath filesystem path to public webpack dir
	FsPath string
	// WebPath http path to public webpack dir
	WebPath string
	// Plugin webpack plugin to use, can be stats or manifest
	Plugin string
	// IgnoreMissing ignore assets missing on manifest or fail on them
	IgnoreMissing bool
	// Verbose - show more info
	Verbose bool
	// IsDev - true to use webpack-serve or webpack-dev-server, false to use filesystem and manifest.json
	IsDev bool

	initDone        bool
	preloadedAssets map[string][]string
}

var AssetHelper func(string) (template.HTML, error)

// Init Set current environment and preload manifest
func Init(dev bool) {
	if Plugin == "deprecated-stats" {
		Plugin = "stats"
		log.Println("go-webpack: default plugin will be changed to manifest instead of stats-plugin")
		log.Println("go-webpack: to continue using stats-plugin, please set webpack.Plugin = 'stats' explicitly")
	}
	isDev = dev

	AssetHelper = GetAssetHelper(&Config{
		DevHost:       DevHost,
		FsPath:        FsPath,
		WebPath:       WebPath,
		Plugin:        Plugin,
		IgnoreMissing: IgnoreMissing,
		Verbose:       Verbose,
		IsDev:         dev,
	})
}

func BasicConfig(host, path, webPath string) *Config {
	return &Config{
		DevHost:       host,
		FsPath:        path,
		WebPath:       webPath,
		Plugin:        "manifest",
		IgnoreMissing: true,
		Verbose:       true,
		IsDev:         isDev,
	}
}

// AssetHelper renders asset tag with url from webpack manifest to the page

func readManifest(conf *Config) (map[string][]string, error) {
	return reader.Read(conf.Plugin, conf.DevHost, conf.FsPath, conf.WebPath, conf.IsDev)
}

func GetAssetHelper(conf *Config) func(string) (template.HTML, error) {
	var err error
	if conf.IsDev {
		// Try to preload manifest, so we can show an error if webpack-dev-server is not running
		_, err = readManifest(conf)
	} else {
		conf.preloadedAssets, err = readManifest(conf)
	}
	if err != nil {
		log.Println(err)
	}
	initDone = true

	return func(key string) (template.HTML, error) {
		var err error

		if !initDone {
			return "", errors.New("Please call webpack.Init() first (see readme)")
		}

		var assets map[string][]string
		if isDev {
			assets, err = readManifest(conf)
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
				log.Printf("%s. Manifest contents:", message)
				for k, a := range assets {
					log.Printf("%s: %s", k, a)
				}
			}
			if IgnoreMissing {
				return template.HTML(""), nil
			}
			return template.HTML(""), errors.New(message)
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
}
