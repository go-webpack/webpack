package manifest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidManifest(t *testing.T) {
	validManifest := `
	{
		"site.js": "site.12345.js",
		"site2.js": "site2.12346.js"
	}
	`

	_, err := unmarshalManifest([]byte(validManifest))
	if err != nil {
		t.Fatalf("shouldn't fail to unmarshal a valid manifest json")
	}
}

func TestReadForEnvironmentDevUsesHTTPManifest(t *testing.T) {
	tmp := t.TempDir()
	writeManifest(t, tmp, `{"site.js":"site.old.js"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webpack/manifest.json" {
			t.Fatalf("unexpected manifest path: %s", r.URL.Path)
		}
		w.Write([]byte(`{"site.js":"site.new.js"}`))
	}))
	defer server.Close()

	assets, err := ReadForEnvironment(true, strings.TrimPrefix(server.URL, "http://"), tmp, "webpack")
	if err != nil {
		t.Fatalf("should read dev manifest from HTTP: %v", err)
	}

	got := assets["site.js"]
	if len(got) != 1 || got[0] != "site.new.js" {
		t.Fatalf("expected HTTP manifest asset, got %#v", got)
	}
}

func TestReadForEnvironmentProdUsesFilesystemManifest(t *testing.T) {
	tmp := t.TempDir()
	writeManifest(t, tmp, `{"site.js":"site.prod.js"}`)

	assets, err := ReadForEnvironment(false, "127.0.0.1:1", tmp, "webpack")
	if err != nil {
		t.Fatalf("should read prod manifest from filesystem: %v", err)
	}

	got := assets["site.js"]
	if len(got) != 1 || got[0] != "site.prod.js" {
		t.Fatalf("expected filesystem manifest asset, got %#v", got)
	}
}

func writeManifest(t *testing.T, dir, contents string) {
	t.Helper()
	if err := ioutil.WriteFile(filepath.Join(dir, "manifest.json"), []byte(contents), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
}
