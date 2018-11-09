package manifest

import "testing"

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
