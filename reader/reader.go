package reader

import (
	"errors"

	"github.com/go-webpack/webpack/reader/manifest"
	"github.com/go-webpack/webpack/reader/stats"
)

// Read webpack asset manifest
func Read(plugin, host, fsPath, webPath string, isDev bool) (map[string][]string, error) {
	//log.Println("read", plugin, isDev)
	if plugin == "stats" {
		return stats.Read(isDev, host, fsPath, webPath)
	} else if plugin == "manifest" {
		return manifest.ReadForEnvironment(isDev, host, fsPath, webPath)
	} else {
		return map[string][]string{}, errors.New("go-webpack: bad plugin type")
	}

}
