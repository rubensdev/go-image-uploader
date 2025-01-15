//go:build dev
// +build dev

package main

import (
	"embed"
	"net/http"
)

var AssetsFS embed.FS

var ManifestJSONStr string

func GetAssetsHandler() http.Handler {
	return nil
}
