package webpack

import (
	"errors"
	"html/template"
	"log"
	"strings"

	"github.com/go-webpack/webpack/helper"
	"github.com/go-webpack/webpack/reader"
)

// AssetList represents a mapping from asset name to asset filenames (maybe multiple assets, with a hash etc)
type AssetList map[string][]string

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

// Config is an instance of go-webpack configuration to use with multiple manifest files (multiple webpack configs)
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
}

// AssetHelper renders asset tag with url from webpack manifest to the page. This is a default assethelper, exported at package level. You can also get your own AssetHelper with webpack.GetAssetHelper
var AssetHelper func(string) (template.HTML, error)

// Init Set current environment and preload manifest
func Init(dev bool) {
	if Plugin == "deprecated-stats" {
		Plugin = "stats"
		log.Println("go-webpack: default plugin will be changed to manifest instead of stats-plugin")
		log.Println("go-webpack: to continue using stats-plugin, please set webpack.Plugin = 'stats' explicitly")
	}

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

// BasicConfig returns a config with basic options set to defaults
func BasicConfig(host, path, webPath string) *Config {
	return &Config{
		DevHost:       host,
		FsPath:        path,
		WebPath:       webPath,
		Plugin:        "manifest",
		IgnoreMissing: true,
		Verbose:       true,
		IsDev:         false,
	}
}

func readManifest(conf *Config) (AssetList, error) {
	//if conf.Verbose {
	//log.Println("go-webpack: reading manifest. Plugin:", conf.Plugin, "dev:", conf.IsDev, "dev host:", conf.DevHost, "fs path:", conf.FsPath, "web path:", conf.WebPath)
	//}
	return reader.Read(conf.Plugin, conf.DevHost, conf.FsPath, conf.WebPath, conf.IsDev)
}

// ErrorFunction returns a template function that returns a fixed error message
func ErrorFunction(err error) func(string) (template.HTML, error) {
	log.Println("go-webpack: error:", err)
	return func(string) (template.HTML, error) {
		return template.HTML(""), err
	}
}

// GetAssetHelper returns an asset helper function based on your config, for use with multiple webpack manifests
func GetAssetHelper(conf *Config) func(string) (template.HTML, error) {
	preloadedAssets := AssetList{}

	var err error
	if conf.IsDev {
		// Try to preload manifest, so we can show an error if webpack-dev-server is not running
		_, err = readManifest(conf)
		if err != nil {
			log.Println(err)
		}
		//if err != nil {
		//return ErrorFunction(err)
		//}
	} else {
		preloadedAssets, err = readManifest(conf)
		// we won't ever re-check assets in this case.  this should be a hard error.
		if err != nil {
			return ErrorFunction(err)
		}
	}

	return createAssetHelper(conf, preloadedAssets)
}

func displayAssetError(conf *Config, key string, assets AssetList) (template.HTML, error) {
	message := "go-webpack: Asset file '" + key + "' not found in manifest"
	if conf.Verbose {
		log.Printf("%s. Manifest contents:", message)
		for k, a := range assets {
			log.Printf("%s: %s", k, a)
		}
	}
	if conf.IgnoreMissing {
		return template.HTML(""), nil
	}
	return template.HTML(""), errors.New(message)
}

func createAssetHelper(conf *Config, preloadedAssets AssetList) func(string) (template.HTML, error) {
	return func(key string) (template.HTML, error) {
		var err error

		var assets AssetList
		if conf.IsDev {
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
			return displayAssetError(conf, key, assets)
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
