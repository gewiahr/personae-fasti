package service

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/gen2brain/webp"
)

func TestProcessUploadedImage(t *testing.T) {
	t.Parallel()

	source := image.NewNRGBA(image.Rect(0, 0, 800, 400))
	for y := 0; y < source.Bounds().Dy(); y++ {
		for x := 0; x < source.Bounds().Dx(); x++ {
			source.SetNRGBA(x, y, color.NRGBA{R: 50, G: 100, B: 150, A: 255})
		}
	}
	var input bytes.Buffer
	if err := png.Encode(&input, source); err != nil {
		t.Fatal(err)
	}

	processed, err := processUploadedImage(bytes.NewReader(input.Bytes()), int64(input.Len()))
	if err != nil {
		t.Fatalf("processUploadedImage() error = %v", err)
	}
	if processed.width != 800 || processed.height != 400 {
		t.Fatalf("full dimensions = %dx%d", processed.width, processed.height)
	}
	thumbnail, err := webp.DecodeConfig(bytes.NewReader(processed.thumbnail))
	if err != nil {
		t.Fatalf("decode thumbnail: %v", err)
	}
	if thumbnail.Width != 512 || thumbnail.Height != 256 {
		t.Fatalf("thumbnail dimensions = %dx%d", thumbnail.Width, thumbnail.Height)
	}
	if processed.checksum == "" || len(processed.full) == 0 || len(processed.thumbnail) == 0 {
		t.Fatal("processed variants or checksum are empty")
	}
}

func TestProcessUploadedImageRejectsOversizedDimensions(t *testing.T) {
	t.Parallel()

	source := image.NewNRGBA(image.Rect(0, 0, maxImageEdge+1, 1))
	var input bytes.Buffer
	if err := png.Encode(&input, source); err != nil {
		t.Fatal(err)
	}
	if _, err := processUploadedImage(bytes.NewReader(input.Bytes()), int64(input.Len())); err == nil {
		t.Fatal("expected oversized dimensions to be rejected")
	}
}
