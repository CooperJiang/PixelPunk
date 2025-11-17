package decode

import (
	"bytes"
	"image"
	"io"
	"strings"

	"pixelpunk/pkg/imagex/iox"

	// Import image format decoders to register them
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/Kodeworks/golang-image-ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// DetectFormat 读取图像尺寸与格式（尽量少读，返回 format 小写）
func DetectFormat(r io.Reader) (width int, height int, format string, err error) {
	data, err := iox.ReadAllWithLimit(r, iox.DefaultMaxReadBytes)
	if err != nil {
		return 0, 0, "", err
	}
	cfg, fmtName, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, "", err
	}
	return cfg.Width, cfg.Height, strings.ToLower(fmtName), nil
}
