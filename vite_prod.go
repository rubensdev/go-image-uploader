//go:build !dev
// +build !dev

package main

import (
	"embed"
)

//go:embed dist/assets/*
var AssetsFS embed.FS

//go:embed dist/.vite/manifest.json
var ManifestJSONStr string
