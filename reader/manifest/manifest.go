package manifest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

// Read webpack-manifest-plugin format manifest
func Read(path string) (map[string][]string, error) {
	log.Println("read:", path+"/manifest.json")
	data, err := ioutil.ReadFile(path + "/manifest.json")
	if err != nil {
		return nil, errors.Wrap(err, "go-webpack: Error when loading manifest from file")
	}

	return unmarshalManifest(data)
}

func unmarshalManifest(data []byte) (map[string][]string, error) {
	response := make(map[string]string, 0)
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, errors.Wrap(err, "go-webpack: Error unmarshaling manifest file")
	}

	assets := make(map[string][]string, len(response))
	for key, value := range response {
		//log.Println("found asset", key, value)
		if !strings.HasSuffix(value, ".map") {
			assets[key] = []string{value}
		}
	}
	return assets, nil
}
