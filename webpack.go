package wp

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

type wpResponse struct {
	Errors            []string                    `json:"errors"`
	Warning           []string                    `json:"warnings"`
	Version           string                      `json:"version"`
	Hash              string                      `json:"hash"`
	PublicPath        string                      `json:"publicPath"`
	AssetsByChunkName map[string]*json.RawMessage `json:"assetsByChunkName"`
	Assets            []*json.RawMessage          `json:"assets"`
}

const host = "localhost:3808"

var c *fasthttp.HostClient
var dev bool
var assets map[string][]string
var webpackBase string

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func devManifest() (data []byte) {
	manifestUrl := fmt.Sprint("http://", host, "/webpack/manifest.json")
	statusCode, body, err := c.Get(nil, manifestUrl)
	if err != nil {
		log.Fatalf("Error when loading manifest from %s: %s", manifestUrl, err)
	}
	if statusCode != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
	}
	return body
}

func prodManifest() (data []byte) {
	body, err := ioutil.ReadFile("./public/webpack/manifest.json")
	if err != nil {
		log.Fatalf("Error when loading manifest from file: %s", err)
	}
	return body
}

func Manifest() map[string][]string {
	var data []byte
	if dev {
		data = devManifest()
	} else {
		data = prodManifest()
	}
	resp := wpResponse{}
	json.Unmarshal(data, &resp)
	webpackBase = resp.PublicPath

	ast := make(map[string][]string, len(resp.AssetsByChunkName))
	var err error
	for akey, aval := range resp.AssetsByChunkName {
		var d []string
		err = json.Unmarshal(*aval, &d)
		if err != nil {
			//log.Fatalf("Error when parsing manifest for %s: %s %s", akey, err, aval)
			//continue
			var sd string
			err = json.Unmarshal(*aval, &sd)
			if err != nil {
				log.Fatalf("Error when parsing manifest for %s: %s %s", akey, err, aval)
				continue
			}
			d = []string{sd}
		}
		ast[akey] = Filter(d, func(v string) bool {
			return !strings.Contains(v, ".map")
		})
		//ast[akey] = d
	}
	//log.Println(ast)
	return ast
}

func AssetHelper(key string) (template.HTML, error) {
	var ast map[string][]string
	if dev {
		ast = Manifest()
	} else {
		ast = assets
	}

	dat := strings.Split(key, ".")

	buf := []string{}
	var err error
	v, ok := ast[dat[0]]
	if !ok {
		return "", errors.New(fmt.Sprint("asset file ", dat[0], " not found in manifest"))
	}
	for _, s := range v {
		if dat[1] == "css" {
			if strings.HasSuffix(s, ".css") {
				buf = append(buf, fmt.Sprint("<link type=\"text/css\" rel=\"stylesheet\" href=\"", webpackBase, html.EscapeString(s), "\"></script>"))
			}
		} else if dat[1] == "js" {
			if strings.HasSuffix(s, ".js") {
				buf = append(buf, fmt.Sprint("<script type=\"text/javascript\" src=\"", webpackBase, html.EscapeString(s), "\"></script>"))
			}
		}
	}

	return template.HTML(strings.Join(buf, "\n")), err
}

func Init(is_dev bool) {
	dev = is_dev
	if dev {
		c = &fasthttp.HostClient{
			Addr: host,
		}
		Manifest()
	} else {
		assets = Manifest()
	}
}
