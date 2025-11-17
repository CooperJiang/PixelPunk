package compress

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

// Options 压缩选项
type Options struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
	Preserve  bool
}

// Result 压缩结果
type Result struct {
	Reader io.Reader
	Width  int
	Height int
	Format string
}

// CompressFile 基于尺寸与质量压缩（保持格式）
func CompressFile(input io.Reader, options *Options) (*Result, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}
	file, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	ow, oh := file.Bounds().Dx(), file.Bounds().Dy()
	tw, th := ow, oh
	if options != nil {
		if options.MaxWidth > 0 || options.MaxHeight > 0 {
			if options.Preserve || true {
				if options.MaxWidth > 0 && (ow > options.MaxWidth) {
					r := float64(options.MaxWidth) / float64(ow)
					tw = options.MaxWidth
					th = int(float64(oh) * r)
				}
				if options.MaxHeight > 0 && th > options.MaxHeight {
					r := float64(options.MaxHeight) / float64(th)
					th = options.MaxHeight
					tw = int(float64(tw) * r)
				}
			}
		}
	}
	out := imaging.Resize(file, tw, th, imaging.Lanczos)
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		if options != nil && options.Quality > 0 {
			if err := imaging.Encode(&buf, out, imaging.JPEG, imaging.JPEGQuality(options.Quality)); err != nil {
				return nil, err
			}
		} else {
			if err := imaging.Encode(&buf, out, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
				return nil, err
			}
		}
	case "png":
		if err := imaging.Encode(&buf, out, imaging.PNG); err != nil {
			return nil, err
		}
	default:
		if err := imaging.Encode(&buf, out, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
			return nil, err
		}
		format = "jpeg"
	}
	b := buf.Bytes()
	return &Result{Reader: bytes.NewReader(b), Width: tw, Height: th, Format: format}, nil
}

// CompressToTargetSize 目标大小压缩（简单迭代质量）
func CompressToTargetSize(reader io.Reader, targetSizeMB float64, options *Options) (*Result, error) {
	if targetSizeMB <= 0 {
		return nil, fmt.Errorf("invalid target size")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	quality := 85
	if options != nil && options.Quality > 0 {
		quality = options.Quality
	}
	for q := quality; q >= 40; q -= 5 {
		res, err := CompressFile(bytes.NewReader(data), &Options{MaxWidth: options.MaxWidth, MaxHeight: options.MaxHeight, Quality: q, Preserve: true})
		if err != nil {
			return nil, err
		}
		buf, _ := io.ReadAll(res.Reader)
		if float64(len(buf)) <= targetSizeMB*1024*1024 {
			return &Result{Reader: bytes.NewReader(buf), Width: res.Width, Height: res.Height, Format: res.Format}, nil
		}
	}
	res, err := CompressFile(bytes.NewReader(data), &Options{MaxWidth: options.MaxWidth, MaxHeight: options.MaxHeight, Quality: 60, Preserve: true})
	if err != nil {
		return nil, err
	}
	return res, nil
}
