package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

type ImageValidator struct {
	MaxWidthPixels   int
	MaxHeightPixels  int
	MaxFileSizeMB    float64
	AllowedMimeTypes []string
	CheckWidth       bool
}

func (v *ImageValidator) MaxFileSizeInBytes() int64 {
	return (int64)(math.Floor(v.MaxFileSizeMB * 1024 * 1024))
}

func (v *ImageValidator) ExceedsMaxFileSize(size int64) bool {
	return size > v.MaxFileSizeInBytes()
}

func (v *ImageValidator) ValidateImage(file io.Reader, contentType string) error {
	fmt.Println(contentType, v.AllowedMimeTypes)
	mimeTypeAllowed := false
	for _, allowedType := range v.AllowedMimeTypes {
		if contentType == allowedType {
			mimeTypeAllowed = true
			break
		}
	}

	if !mimeTypeAllowed {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}

	if !v.CheckWidth {
		return nil
	}

	// Decode the image to check dimension
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	if width > v.MaxWidthPixels {
		return fmt.Errorf("image width %d pixels exceeds maximum of %d pixels", width, v.MaxWidthPixels)
	}

	if height > v.MaxHeightPixels {
		return fmt.Errorf("image height %d pixels exceeds maximum of %d pixels", height, v.MaxHeightPixels)
	}

	return nil
}
