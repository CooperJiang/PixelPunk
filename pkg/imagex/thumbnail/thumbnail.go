package thumbnail

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/Kodeworks/golang-image-ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/disintegration/imaging"
	oksvg "github.com/srwiley/oksvg"
	rasterx "github.com/srwiley/rasterx"
)

// Options 缩略图参数
type Options struct {
	Width    int
	Height   int
	Quality  int
	Crop     bool
	Preserve bool
	Format   string // 指定输出，空为自动
}

// Result 缩略图结果
type Result struct {
	Reader io.Reader
	Width  int
	Height int
	Size   int64
	Format string
}

// Generate 生成缩略图（封装 SVG 与位图路径）
func Generate(input []byte, opts Options) (*Result, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	// 优先尝试 SVG
	if looksLikeSVG(input) {
		return renderSVG(input, opts)
	}

	// ICO 特殊处理：尝试提取内部的PNG数据
	if isICOFormat(input) {
		if pngData := extractPNGFromICO(input); pngData != nil {
			// 使用提取的PNG数据进行处理
			file, _, err := image.Decode(bytes.NewReader(pngData))
			if err == nil {
				// ICO文件强制使用PNG格式保留透明度
				return resizeAndEncodeWithFormat(file, opts, "png")
			}
		}
	}

	file, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	return resizeAndEncode(file, opts)
}

func looksLikeSVG(data []byte) bool {
	lower := bytes.ToLower(data)
	if bytes.Contains(lower, []byte("<svg")) {
		return true
	}
	if bytes.HasPrefix(lower, []byte("<?xml")) && bytes.Contains(lower, []byte("<svg")) {
		return true
	}
	return false
}

// isICOFormat 检查是否为ICO格式
func isICOFormat(data []byte) bool {
	// ICO文件头：00 00 01 00 (小端序)
	if len(data) < 4 {
		return false
	}
	return data[0] == 0x00 && data[1] == 0x00 && data[2] == 0x01 && data[3] == 0x00
}

// extractPNGFromICO 从ICO文件中提取PNG数据
func extractPNGFromICO(data []byte) []byte {
	if len(data) < 22 {
		return nil
	}

	// ICO格式分析：
	// 0-1: 保留字段 (00 00)
	// 2-3: 图像类型 (01 00 表示ICO)
	// 4-5: 图像数量
	imageCount := int(data[4]) | int(data[5])<<8
	if imageCount == 0 {
		return nil
	}

	// 每个图像目录条目16字节，从偏移6开始
	for i := 0; i < imageCount; i++ {
		dirOffset := 6 + i*16
		if dirOffset+16 > len(data) {
			break
		}

		// 读取图像数据的偏移量和大小
		dataSize := int(data[dirOffset+8]) | int(data[dirOffset+9])<<8 | int(data[dirOffset+10])<<16 | int(data[dirOffset+11])<<24
		dataOffset := int(data[dirOffset+12]) | int(data[dirOffset+13])<<8 | int(data[dirOffset+14])<<16 | int(data[dirOffset+15])<<24

		if dataOffset+dataSize > len(data) {
			continue
		}

		imageData := data[dataOffset : dataOffset+dataSize]

		// 检查是否为PNG格式 (89 50 4E 47)
		if len(imageData) >= 4 && imageData[0] == 0x89 && imageData[1] == 0x50 && imageData[2] == 0x4E && imageData[3] == 0x47 {
			return imageData
		}
	}

	return nil
}

func renderSVG(input []byte, opts Options) (*Result, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(input))
	if err != nil || icon == nil {
		return nil, fmt.Errorf("decode svg: %w", err)
	}
	vb := icon.ViewBox
	ow, oh := int(vb.W), int(vb.H)
	if ow <= 0 || oh <= 0 {
		ow, oh = 512, 512
	}
	tw, th := computeSize(ow, oh, opts)
	if tw <= 0 {
		tw = 200
	}
	if th <= 0 {
		th = 200
	}
	rgba := image.NewRGBA(image.Rect(0, 0, tw, th))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	icon.SetTarget(0, 0, float64(tw), float64(th))
	scanner := rasterx.NewScannerGV(tw, th, rgba, rgba.Bounds())
	dasher := rasterx.NewDasher(tw, th, scanner)
	icon.Draw(dasher, 1.0)
	return encode(rgba, "png", opts)
}

func resizeAndEncode(file image.Image, opts Options) (*Result, error) {
	ow, oh := file.Bounds().Dx(), file.Bounds().Dy()
	w, h := computeSize(ow, oh, opts)
	var out image.Image
	if opts.Crop {
		out = imaging.Fill(file, w, h, imaging.Center, imaging.Lanczos)
	} else {
		out = imaging.Resize(file, w, h, imaging.Lanczos)
	}
	// 按需选择格式：显式指定优先；否则含透明优先 PNG，其余 JPEG
	format := opts.Format
	if format == "" {
		if imageHasAlpha(out) {
			format = "png"
		} else {
			format = "jpeg"
		}
	}
	return encode(out, format, opts)
}

func resizeAndEncodeWithFormat(file image.Image, opts Options, forceFormat string) (*Result, error) {
	ow, oh := file.Bounds().Dx(), file.Bounds().Dy()
	tw, th := computeSize(ow, oh, opts)
	var out image.Image
	if opts.Crop {
		out = imaging.Fill(file, tw, th, imaging.Center, imaging.Lanczos)
	} else {
		out = imaging.Resize(file, tw, th, imaging.Lanczos)
	}
	return encode(out, forceFormat, opts)
}

func computeSize(ow, oh int, opts Options) (int, int) {
	if opts.Width > 0 && opts.Height > 0 {
		if opts.Preserve {
			wr := float64(opts.Width) / float64(ow)
			hr := float64(opts.Height) / float64(oh)
			if wr < hr {
				return opts.Width, int(float64(oh) * wr)
			}
			return int(float64(ow) * hr), opts.Height
		}
		return opts.Width, opts.Height
	}
	if opts.Width > 0 {
		if opts.Preserve {
			r := float64(opts.Width) / float64(ow)
			return opts.Width, int(float64(oh) * r)
		}
		return opts.Width, opts.Width
	}
	if opts.Height > 0 {
		if opts.Preserve {
			r := float64(opts.Height) / float64(oh)
			return int(float64(ow) * r), opts.Height
		}
		return opts.Height, opts.Height
	}
	size := 200
	if opts.Crop {
		return size, size
	}
	if ow > oh {
		r := float64(size) / float64(ow)
		return size, int(float64(oh) * r)
	}
	r := float64(size) / float64(oh)
	return int(float64(ow) * r), size
}

func encode(file image.Image, format string, opts Options) (*Result, error) {
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		if err := imaging.Encode(&buf, file, imaging.JPEG, imaging.JPEGQuality(safeQuality(opts.Quality))); err != nil {
			return nil, err
		}
	case "png":
		if err := imaging.Encode(&buf, file, imaging.PNG); err != nil {
			return nil, err
		}
	default:
		if err := imaging.Encode(&buf, file, imaging.JPEG, imaging.JPEGQuality(safeQuality(opts.Quality))); err != nil {
			return nil, err
		}
		format = "jpeg"
	}
	b := buf.Bytes()
	return &Result{Reader: bytes.NewReader(b), Width: file.Bounds().Dx(), Height: file.Bounds().Dy(), Size: int64(len(b)), Format: format}, nil
}

func safeQuality(q int) int {
	if q <= 0 {
		return 80
	}
	if q > 100 {
		return 100
	}
	return q
}

func imageHasAlpha(file image.Image) bool {
	b := file.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := file.At(x, y).RGBA()
			if a != 0xffff {
				return true
			}
		}
	}
	return false
}
