package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
)

type config struct {
	port             int
	maxWidthPixels   int
	maxHeightPixels  int
	maxFileSizeMB    float64
	allowedMimeTypes []string
	checkWidth       bool
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "App server port")
	flag.IntVar(&cfg.maxWidthPixels, "max-width-pixels", 1920, "maximum width allowed in pixels")
	flag.IntVar(&cfg.maxHeightPixels, "max-height-pixels", 1080, "maximum height allowd in pixels")
	flag.Float64Var(&cfg.maxFileSizeMB, "max-file-size", 5.00, "maximum file size allowed in Megabytes")

	flag.Func("allowed-mimetypes", "Allowed mimetypes", func(val string) error {
		if val != "" {
			cfg.allowedMimeTypes = strings.Fields(val)
		}
		return nil
	})

	flag.Parse()

	if len(cfg.allowedMimeTypes) == 0 {
		cfg.allowedMimeTypes = []string{"image/jpeg", "image/png", "image/gif"}
	}

	validator := &ImageValidator{
		MaxWidthPixels:   cfg.maxWidthPixels,
		MaxHeightPixels:  cfg.maxHeightPixels,
		MaxFileSizeMB:    cfg.maxFileSizeMB,
		AllowedMimeTypes: cfg.allowedMimeTypes,
		CheckWidth:       cfg.checkWidth,
	}

	router := http.NewServeMux()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	component := Hello("Rubens")

	fs := http.FileServer(http.Dir("./static"))

	imgHandler := NewImageHandler(logger, validator)
	router.Handle("/", templ.Handler(component))
	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/upload", imgHandler.Handle)
	router.HandleFunc("/upload/multiple", imgHandler.HandleMultiple)

	logger.Info("Server starting", "port", cfg.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), router))
}
