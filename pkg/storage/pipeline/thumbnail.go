package pipeline

import (
	"io"

	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/imagex/convert"
	"pixelpunk/pkg/imagex/thumbnail"
)

// Options 定义缩略图生成参数
type Options struct {
	Width           int
	Height          int
	Quality         int
	EnableWebP      bool
	FallbackOnError bool
}

// Result 缩略图生成结果
type Result struct {
	Data          []byte // 缩略图数据
	Format        string // 格式
	Failed        bool   // 是否失败
	FailureReason string // 失败原因
}

// GenerateOrFallback 生成缩略图；失败时返回内置失败图（webp）
func GenerateOrFallback(input []byte, opts Options) ([]byte, string, error) {
	q := opts.Quality
	if q <= 0 {
		q = 85
	}
	w := opts.Width
	if w <= 0 {
		w = 1200
	}
	h := opts.Height
	if h <= 0 {
		h = 900
	}

	thumbRes, err := thumbnail.Generate(input, thumbnail.Options{
		Width:    w,
		Height:   h,
		Quality:  q,
		Crop:     false,
		Preserve: true,
	})
	if err != nil {
		if opts.FallbackOnError {
			failData, ferr := assets.GetDefaultFileData(assets.FileTypeFail)
			if ferr == nil {
				return failData, "webp", nil
			}
		}
		return nil, "", err
	}

	thumbBytes, _ := io.ReadAll(thumbRes.Reader)
	format := thumbRes.Format
	if format == "" {
		format = "jpg"
	}

	if opts.EnableWebP {
		if webp, err := convert.ToWebP(thumbBytes, convert.WebPOptions{Quality: q}); err == nil && webp.Converted {
			buf, _ := io.ReadAll(webp.Reader)
			thumbBytes = buf
			format = "webp"
		}
	}

	return thumbBytes, format, nil
}

// GenerateWithResult 生成缩略图并返回详细结果（包含失败信息）
func GenerateWithResult(input []byte, opts Options) *Result {
	q := opts.Quality
	if q <= 0 {
		q = 85
	}
	w := opts.Width
	if w <= 0 {
		w = 1200
	}
	h := opts.Height
	if h <= 0 {
		h = 900
	}

	thumbRes, err := thumbnail.Generate(input, thumbnail.Options{
		Width:    w,
		Height:   h,
		Quality:  q,
		Crop:     false,
		Preserve: true,
	})
	if err != nil {
		return &Result{
			Failed:        true,
			FailureReason: err.Error(),
		}
	}

	thumbBytes, _ := io.ReadAll(thumbRes.Reader)
	format := thumbRes.Format
	if format == "" {
		format = "jpg"
	}

	if opts.EnableWebP {
		if webp, err := convert.ToWebP(thumbBytes, convert.WebPOptions{Quality: q}); err == nil && webp.Converted {
			buf, _ := io.ReadAll(webp.Reader)
			thumbBytes = buf
			format = "webp"
		}
	}

	return &Result{
		Data:   thumbBytes,
		Format: format,
		Failed: false,
	}
}
