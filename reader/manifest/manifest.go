package manifest

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type assetList map[string][]string
type assetResponse map[string]string

// Read webpack-manifest-plugin format manifest
func Read(path string) (assetList, error) {
	//log.Println("read:", path+"/manifest.json")
	data, err := ioutil.ReadFile(path + "/manifest.json")
	if err != nil {
		return nil, errors.Wrap(err, "go-webpack: Error when loading manifest from file")
	}

	return unmarshalManifest(data)
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
