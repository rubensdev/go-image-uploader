//go:build !dev
// +build !dev

package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/assets/*
var AssetsFS embed.FS

//go:embed dist/.vite/manifest.json
var ManifestJSONStr string

func GetAssetsHandler() http.Handler {
	stripped, err := fs.Sub(AssetsFS, "dist/assets")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.FS(stripped))
	return http.StripPrefix("/assets/", fs)
}
