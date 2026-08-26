package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"github.com/gen2brain/webp"
	"github.com/rwcarlsen/goexif/exif"
	xdraw "golang.org/x/image/draw"
)

const (
	maxImageEdge   = 12_000
	maxImagePixels = 40_000_000
	thumbnailEdge  = 512
)

type processedImage struct {
	full      []byte
	thumbnail []byte
	width     int
	height    int
	checksum  string
}

func processUploadedImage(reader io.Reader, maxFileBytes int64) (*processedImage, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maxFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	if int64(len(data)) > maxFileBytes {
		return nil, fmt.Errorf("file exceeds %d bytes", maxFileBytes)
	}

	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return nil, fmt.Errorf("unsupported image format")
	}
	if contentType == "image/webp" && animatedWebP(data) {
		return nil, fmt.Errorf("animated WebP is not supported")
	}

	var config image.Config
	if contentType == "image/webp" {
		config, err = webp.DecodeConfig(bytes.NewReader(data))
	} else {
		config, _, err = image.DecodeConfig(bytes.NewReader(data))
	}
	if err != nil || config.Width <= 0 || config.Height <= 0 {
		return nil, fmt.Errorf("invalid image")
	}
	if config.Width > maxImageEdge || config.Height > maxImageEdge || int64(config.Width)*int64(config.Height) > maxImagePixels {
		return nil, fmt.Errorf("image dimensions are too large")
	}

	var decoded image.Image
	if contentType == "image/webp" {
		decoded, err = webp.Decode(bytes.NewReader(data), webp.Options{AutoRotate: true})
	} else {
		decoded, _, err = image.Decode(bytes.NewReader(data))
		if err == nil && contentType == "image/jpeg" {
			decoded = applyOrientation(decoded, jpegOrientation(data))
		}
	}
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	var full bytes.Buffer
	if err := webp.Encode(&full, decoded, webp.Options{Quality: 82, Method: 4}); err != nil {
		return nil, fmt.Errorf("encode image: %w", err)
	}
	thumbnail := resizeToEdge(decoded, thumbnailEdge)
	var thumb bytes.Buffer
	if err := webp.Encode(&thumb, thumbnail, webp.Options{Quality: 78, Method: 4}); err != nil {
		return nil, fmt.Errorf("encode thumbnail: %w", err)
	}

	checksum := sha256.New()
	_, _ = checksum.Write(full.Bytes())
	_, _ = checksum.Write(thumb.Bytes())
	bounds := decoded.Bounds()
	return &processedImage{
		full:      full.Bytes(),
		thumbnail: thumb.Bytes(),
		width:     bounds.Dx(),
		height:    bounds.Dy(),
		checksum:  hex.EncodeToString(checksum.Sum(nil)),
	}, nil
}

func resizeToEdge(source image.Image, maxEdge int) image.Image {
	bounds := source.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width <= maxEdge && height <= maxEdge {
		return source
	}
	newWidth, newHeight := width, height
	if width >= height {
		newWidth = maxEdge
		newHeight = max(1, height*maxEdge/width)
	} else {
		newHeight = maxEdge
		newWidth = max(1, width*maxEdge/height)
	}
	destination := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xdraw.CatmullRom.Scale(destination, destination.Bounds(), source, bounds, draw.Src, nil)
	return destination
}

func jpegOrientation(data []byte) int {
	metadata, err := exif.Decode(bytes.NewReader(data))
	if err != nil || metadata == nil {
		return 1
	}
	tag, err := metadata.Get(exif.Orientation)
	if err != nil {
		return 1
	}
	orientation, err := tag.Int(0)
	if err != nil || orientation < 1 || orientation > 8 {
		return 1
	}
	return orientation
}

func applyOrientation(source image.Image, orientation int) image.Image {
	if orientation == 1 {
		return source
	}
	bounds := source.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	destinationWidth, destinationHeight := width, height
	if orientation >= 5 {
		destinationWidth, destinationHeight = height, width
	}
	destination := image.NewNRGBA(image.Rect(0, 0, destinationWidth, destinationHeight))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx, dy := x, y
			switch orientation {
			case 2:
				dx = width - 1 - x
			case 3:
				dx, dy = width-1-x, height-1-y
			case 4:
				dy = height - 1 - y
			case 5:
				dx, dy = y, x
			case 6:
				dx, dy = height-1-y, x
			case 7:
				dx, dy = height-1-y, width-1-x
			case 8:
				dx, dy = y, width-1-x
			}
			destination.Set(dx, dy, source.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return destination
}

func animatedWebP(data []byte) bool {
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return false
	}
	for offset := 12; offset+8 <= len(data); {
		chunk := string(data[offset : offset+4])
		if chunk == "ANIM" || chunk == "ANMF" {
			return true
		}
		size := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		offset += 8 + size + size%2
	}
	return false
}
