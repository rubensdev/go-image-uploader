package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"rubensdev.com/go-image-processing/templates"
	"rubensdev.com/go-image-processing/templates/vite"
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
	flag.Float64Var(&cfg.maxFileSizeMB, "max-file-size", 5.00, "maximum file size allowed in Megabytes")
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

	vm, err := vite.NewManager(ManifestJSONStr, cfg.env)
	if err != nil {
		log.Fatalf("Error initializing vite manager: %v", err.Error())
	}

	homeVD := templates.ViewData{
		Title: "GOAT Image Uploader",
		Lang:  "es",
		Meta: templates.Metadata{
			"allowed_mimetypes": cfg.allowedMimeTypes,
			"max_file_size":     fmt.Sprintf("%.2f", cfg.maxFileSizeMB),
			"upload_endpoint":   "/upload",
		},
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/favicon.ico" {
			http.NotFound(w, r)
			return
		}

		if vm.InDevMode() {
			vm.SetDevServerURL(strings.Split(r.Host, ":")[0])
		}
		ctx := context.WithValue(context.Background(), templates.ViteManagerCtx, vm)
		templates.Home(homeVD).Render(ctx, w)
	})

	if cfg.env == "production" {
		router.Handle("/assets/", GetAssetsHandler())
	}

	imgHandler := NewImageHandler(logger, validator)
	router.HandleFunc("/upload", imgHandler.Handle)
	//router.HandleFunc("/upload/multiple", imgHandler.HandleMultiple)

	logger.Info("Server starting", "port", cfg.port)

	err = http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.port),
		router,
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
