package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"rubensdev.com/go-image-processing/templates"
	"rubensdev.com/go-image-processing/templates/manifest"
)

type config struct {
	port             int
	maxWidthPixels   int
	maxHeightPixels  int
	maxFileSizeMB    float64
	allowedMimeTypes []string
	checkWidth       bool
	env              string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "App server port")
	flag.IntVar(&cfg.maxWidthPixels, "max-width-pixels", 1920, "maximum width allowed in pixels")
	flag.IntVar(&cfg.maxHeightPixels, "max-height-pixels", 1080, "maximum height allowd in pixels")
	flag.Float64Var(&cfg.maxFileSizeMB, "max-file-size", 1.00, "maximum file size allowed in Megabytes")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

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

	devHost := GetOutboundIP().String()

	mm, err := manifest.NewManager(ManifestJSONStr, cfg.env, devHost)
	if err != nil {
		log.Fatalf("Error loading manifest: %v", err.Error())
	}

	uploadsFs := http.FileServer(http.Dir("./uploads"))

	homeVD := templates.ViewData{
		Title: "Desk Setup",
		Lang:  "es",
		Meta: templates.Metadata{
			"allowed_mimetypes": cfg.allowedMimeTypes,
			"max_file_size":     fmt.Sprintf("%.2f", cfg.maxFileSizeMB),
			"upload_endpoint":   "/upload",
		},
	}
	homeView := templates.Home(homeVD)

	router.Handle("/", templ.Handler(homeView))

	if cfg.env == "production" {
		stripped, err := fs.Sub(AssetsFS, "dist/assets")
		if err != nil {
			log.Fatal(err)
		}
		fs := http.FileServer(http.FS(stripped))
		router.Handle("/assets/", http.StripPrefix("/assets/", fs))
	}

	router.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsFs))

	imgHandler := NewImageHandler(logger, validator)
	router.HandleFunc("/upload", imgHandler.Handle)
	router.HandleFunc("/upload/multiple", imgHandler.HandleMultiple)

	logger.Info("Server starting", "port", cfg.port)

	err = http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.port),
		templates.WithManifestManager(mm)(router),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
