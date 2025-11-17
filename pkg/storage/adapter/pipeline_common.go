package adapter

import (
	"bytes"
	"io"
	"strings"

	"pixelpunk/pkg/imagex/compress"
	"pixelpunk/pkg/imagex/convert"
	"pixelpunk/pkg/imagex/decode"
	"pixelpunk/pkg/storage/pipeline"
)

// processUploadData applies optional compression and WebP conversion, and returns
// processed bytes along with detected width/height/format.
func processUploadData(original []byte, req *UploadRequest) (processed []byte, width int, height int, format string) {
	processed = original

	// 检查并转换 HEIC/HEIF 为 JPEG
	if convert.IsHEICFormat(original) {
		heicResult, err := convert.ToJPEGFromHEIC(original, convert.HEICToJPEGOptions{Quality: 95})
		if err == nil && heicResult.Converted {
			if buf, e := io.ReadAll(heicResult.Reader); e == nil {
				processed = buf
				format = "jpeg"
				width = heicResult.Width
				height = heicResult.Height

				if req != nil {
					req.FileName = replaceHEICExtension(req.FileName)
				}
			}
		}
	}

	// detect initial format (如果还没有检测过)
	if format == "" {
		if w, h, f, err := decode.DetectFormat(bytes.NewReader(processed)); err == nil {
			width, height, format = w, h, f
		}
	}

	// compress if requested
	if req != nil && req.Options != nil && req.Options.Compress {
		if cr, err := compress.CompressToTargetSize(bytes.NewReader(processed), 5.0, &compress.Options{
			MaxWidth: req.Options.MaxWidth, MaxHeight: req.Options.MaxHeight, Quality: req.Options.Quality,
		}); err == nil {
			if buf, e := io.ReadAll(cr.Reader); e == nil {
				processed = buf
				if cr.Width > 0 {
					width = cr.Width
				}
				if cr.Height > 0 {
					height = cr.Height
				}
				if cr.Format != "" {
					format = cr.Format
				}
			}
		}
	}
	// optional webp
	if req != nil && req.Options != nil && req.Options.WebPEnabled {
		if wr, err := convert.ToWebP(processed, convert.WebPOptions{Quality: req.Options.Quality}); err == nil && wr.Converted {
			if buf, e := io.ReadAll(wr.Reader); e == nil {
				processed = buf
				format = "webp"
			}
		}
	}
	return
}

// replaceHEICExtension 将 .heic 或 .HEIF 扩展名替换为 .jpg
func replaceHEICExtension(filename string) string {
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".heic") {
		return filename[:len(filename)-5] + ".jpg"
	}
	if strings.HasSuffix(lower, ".heif") {
		return filename[:len(filename)-5] + ".jpg"
	}
	return filename
}

// buildThumbnailBytes generates a thumbnail with fallback, returning bytes and format.
// The input data should be the best available source (usually original data).
func buildThumbnailBytes(source []byte, req *UploadRequest) (thumbBytes []byte, thumbFormat string) {
	tw, th := 1200, 900
	tq := 85
	if req != nil && req.Options != nil {
		if req.Options.ThumbWidth > 0 {
			tw = req.Options.ThumbWidth
		}
		if req.Options.ThumbHeight > 0 {
			th = req.Options.ThumbHeight
		}
		if req.Options.ThumbQuality > 0 {
			tq = req.Options.ThumbQuality
		}
	}
	tb, tf, _ := pipeline.GenerateOrFallback(source, pipeline.Options{Width: tw, Height: th, Quality: tq, EnableWebP: true, FallbackOnError: true})
	return tb, tf
}
