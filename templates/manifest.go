package templates

import (
	"encoding/json"
	"log"
	"os"
)

type ViteManifest map[string]struct {
	File     string
	CssFiles []string `json:"css"`
	Src      string
}

var IsProd = os.Getenv("ENV") == "production"
var Manifest ViteManifest

func ParseManifest(path string) (ViteManifest, error) {
	binary, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest ViteManifest
	err = json.Unmarshal(binary, &manifest)
	return manifest, err
}

func init() {
	if !IsProd {
		return
	}

	manifest, err := ParseManifest("./dist/.vite/manifest.json")
	if err != nil {
		log.Fatal(err)
	}
	Manifest = manifest
}
