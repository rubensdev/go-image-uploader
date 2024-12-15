package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"

	"golang.org/x/image/draw"
)

type ImageHandler struct {
	logger    *slog.Logger
	validator *ImageValidator
}

func NewImageHandler(l *slog.Logger, v *ImageValidator) *ImageHandler {
	return &ImageHandler{
		logger:    l,
		validator: v,
	}
}

func (h *ImageHandler) Handle(w http.ResponseWriter, r *http.Request) {
	maxFileSizeInBytes := h.validator.MaxFileSizeInBytes()

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSizeInBytes)
	if err := r.ParseMultipartForm(maxFileSizeInBytes); err != nil {
		h.logger.Info(fmt.Sprintf("Error parsing body request %v", err.Error()))
		http.Error(w, "File is too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if h.validator.ExceedsMaxFileSize(header.Size) {
		fileSizeMB := float64(header.Size) / (1024 * 1024)
		http.Error(
			w,
			fmt.Sprintf("File size %.2f MB exceeds maximum of %.2f MB", fileSizeMB, h.validator.MaxFileSizeMB),
			http.StatusBadRequest,
		)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	h.logger.Info("Uploaded image", "name", header.Filename, "mimeType", mimeType)

	if err := h.validator.ValidateImage(file, mimeType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(fmt.Sprintf("uploads/%s", header.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Success().Render(r.Context(), w)
}

func (h *ImageHandler) HandleMultiple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Info(fmt.Sprintf("Error parsing body request %v", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the filehandlers.
	// They are accesible only after ParseMultipartForm is called
	files := r.MultipartForm.File["image"]

	for _, fileHeader := range files {
		if h.validator.ExceedsMaxFileSize(fileHeader.Size) {
			msg := fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than %.2fMB in size", fileHeader.Filename, h.validator.MaxFileSizeMB)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		mimeType := fileHeader.Header.Get("Content-Type")
		h.logger.Info("Uploaded image", "name", fileHeader.Filename, "mimeType", mimeType)

		if err := h.validator.ValidateImage(file, mimeType); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var (
			src image.Image
		)

		switch mimeType {
		case "image/jpeg":
			src, err = jpeg.Decode(file)
		case "image/png":
			src, err = png.Decode(file)
		}

		if err != nil {
			h.logger.Info("Error decoding image", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// THIS IS GOOD FOR CREATING thumbnails
		ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)

		height := int(math.Round(float64(h.validator.MaxWidthPixels) * ratio))

		newImage := image.NewRGBA(image.Rect(0, 0, h.validator.MaxWidthPixels, height))

		// Try with ApproxBiLinear, BiLinear and CatmullRom
		draw.NearestNeighbor.Scale(newImage, newImage.Rect, src, src.Bounds(), draw.Over, nil)

		buf := bytes.NewBuffer(nil)

		switch mimeType {
		case "image/jpeg":
			h.logger.Info("Encoding jpg")
			err = jpeg.Encode(buf, newImage, nil)
		case "image/png":
			h.logger.Info("Encoding png")
			err = png.Encode(buf, newImage)
		}

		if err != nil {
			h.logger.Info("Error encoding image", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(fmt.Sprintf("uploads/%s", fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(dst, io.TeeReader(buf, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successfull"))
}

/**
	- To restrict the type of the uploaded file

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" { {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
**/
