package webpack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestManifestAssetHelper(t *testing.T) {
	assets := map[string][]string{
		"main.js": {"main.1.js", "main.2.js"},
	}

	helper := createAssetHelper(&Config{
		Plugin: "manifest",
	}, assets)

	html, err := helper("main.js")
	if err != nil {
		t.Fatalf("error %v returned from asset helper for valid asset", err)
	}
	expectedHTML :=
		`<script type="text/javascript" src="main.1.js"></script>
<script type="text/javascript" src="main.2.js"></script>`

	if string(html) != expectedHTML {
		t.Fatalf("unexpected <script> tags\nexpected:\n%s\nactual:\n%s", expectedHTML, html)
	}

	// IgnoreMissing = false
	_, err = helper("maiin.js")
	if err == nil {
		t.Fatalf("error nil when it shouldn't have been")
	}
}

func TestManifestAssetHelperDevReloadsHTTPManifestEachCall(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webpack/manifest.json" {
			t.Fatalf("unexpected manifest path: %s", r.URL.Path)
		}
		requests++
		w.Write([]byte(fmt.Sprintf(`{"main.js":"main.%d.js"}`, requests)))
	}))
	defer server.Close()

	helper := GetAssetHelper(&Config{
		DevHost:       strings.TrimPrefix(server.URL, "http://"),
		WebPath:       "webpack",
		Plugin:        "manifest",
		IgnoreMissing: false,
		IsDev:         true,
	})

	first, err := helper("main.js")
	if err != nil {
		t.Fatalf("first helper call failed: %v", err)
	}
	second, err := helper("main.js")
	if err != nil {
		t.Fatalf("second helper call failed: %v", err)
	}

	if !strings.Contains(string(first), "main.2.js") {
		t.Fatalf("expected first helper call to use first reloaded manifest, got %s", first)
	}
	if !strings.Contains(string(second), "main.3.js") {
		t.Fatalf("expected second helper call to use second reloaded manifest, got %s", second)
	}
	if requests != 3 {
		t.Fatalf("expected preload plus two helper reads, got %d requests", requests)
	}
}
