package convert

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type WebPOptions struct {
	Quality          int
	Lossless         bool
	CompressionLevel int
}

type WebPResult struct {
	Reader    io.Reader
	Converted bool
	Size      int64
}

func ToWebP(input []byte, opts WebPOptions) (*WebPResult, error) {
	file, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	var enc *encoder.Options
	if opts.Lossless {
		enc, err = encoder.NewLosslessEncoderOptions(encoder.PresetDefault, clampLevel(opts.CompressionLevel))
	} else {
		enc, err = encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(safeQ(opts.Quality)))
	}
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := webp.Encode(&buf, file, enc); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	return &WebPResult{Reader: bytes.NewReader(b), Converted: true, Size: int64(len(b))}, nil
}

// AsOriginal 当WebP失败时，回退原格式编码
func AsOriginal(file image.Image, format string, quality int) (*WebPResult, error) {
	var buf bytes.Buffer
	switch format {
	case "png":
		if err := imaging.Encode(&buf, file, imaging.PNG); err != nil {
			return nil, err
		}
	default:
		if err := imaging.Encode(&buf, file, imaging.JPEG, imaging.JPEGQuality(safeQ(quality))); err != nil {
			return nil, err
		}
	}
	b := buf.Bytes()
	return &WebPResult{Reader: bytes.NewReader(b), Converted: false, Size: int64(len(b))}, nil
}

func safeQ(q int) int {
	if q <= 0 {
		return 80
	}
	if q > 100 {
		return 100
	}
	return q
}

func clampLevel(l int) int {
	if l < 0 {
		return 0
	}
	if l > 6 {
		return 6
	}
	return l
}
