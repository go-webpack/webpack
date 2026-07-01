package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type assetList map[string][]string
type assetResponse map[string]string

// Read webpack-manifest-plugin format manifest
func Read(path string) (assetList, error) {
	//log.Println("read:", path+"/manifest.json")
	data, err := prodManifest(path)
	if err != nil {
		return nil, err
	}

	return unmarshalManifest(data)
}

// ReadForEnvironment reads manifest from webpack-dev-server in dev mode and from disk in production.
func ReadForEnvironment(isDev bool, host, fsPath, webPath string) (assetList, error) {
	var data []byte
	var err error

	if isDev {
		data, err = devManifest(host, webPath)
	} else {
		data, err = prodManifest(fsPath)
	}
	if err != nil {
		return nil, err
	}

	return unmarshalManifest(data)
}

func devManifest(host, webPath string) ([]byte, error) {
	manifestURL := fmt.Sprint("http://", host, "/", webPath, "/manifest.json")
	resp, err := http.Get(manifestURL)
	if err != nil {
		return []byte{}, errors.Wrap(err, fmt.Sprintf("go-webpack: Error when loading manifest from url %s", manifestURL))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("go-webpack: Unexpected status code: %d. Expecting %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, fmt.Sprintf("go-webpack: Error when loading manifest from url %s", manifestURL))
	}
	return body, nil
}

func prodManifest(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path + "/manifest.json")
	if err != nil {
		return nil, errors.Wrap(err, "go-webpack: Error when loading manifest from file")
	}
	return data, nil
}

func unmarshalManifest(data []byte) (assetList, error) {
	response := make(assetResponse, 0)
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, errors.Wrap(err, "go-webpack: Error unmarshaling manifest file")
	}

	assets := make(assetList, len(response))
	for key, value := range response {
		//log.Println("found asset", key, value)
		if !strings.HasSuffix(value, ".map") {
			assets[key] = []string{value}
		}
	}
	return assets, nil
}
