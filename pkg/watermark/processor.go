package watermark

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

type Processor struct {
	watermarkImagePath  string
	maxImageSize        int
	defaultFontSize     float64
	enableTextWatermark bool // 预留开关（当前不支持文字水印）
}

func NewProcessor() *Processor {
	return &Processor{
		watermarkImagePath: "./assets/watermarks",
		maxImageSize:       8192, // 支持现代相机/手机的8K分辨率
		defaultFontSize:    24,
	}
}

func (p *Processor) SetWatermarkImagePath(path string) {
	p.watermarkImagePath = path
}

func (p *Processor) SetMaxImageSize(size int) {
	p.maxImageSize = size
}

func (p *Processor) SetEnableTextWatermark(enabled bool) {
	p.enableTextWatermark = enabled
}

func (p *Processor) ProcessImage(file image.Image, config *WatermarkConfig) (image.Image, error) {
	if config == nil || !config.Enabled {
		return file, nil
	}

	if err := p.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	startTime := time.Now()
	defer func() {
		processTime := time.Since(startTime).Milliseconds()
		_ = processTime
	}()

	bounds := file.Bounds()
	if bounds.Max.X > p.maxImageSize || bounds.Max.Y > p.maxImageSize {
		return nil, fmt.Errorf("文件尺寸过大，最大支持 %dx%d", p.maxImageSize, p.maxImageSize)
	}

	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, file, bounds.Min, draw.Over)

	if err := p.applyImageWatermark(newImg, config); err != nil {
		return nil, fmt.Errorf("应用文件水印失败: %w", err)
	}

	return newImg, nil
}

func (p *Processor) ProcessImageBytes(imgData []byte, config *WatermarkConfig) ([]byte, error) {
	if config == nil || !config.Enabled {
		return imgData, nil
	}

	file, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("文件解码失败: %w", err)
	}

	processedImg, err := p.ProcessImage(file, config)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		if p.hasTransparency(processedImg) {
			err = png.Encode(&buf, processedImg)
		} else {
			err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: 95})
		}
	case "png":
		err = png.Encode(&buf, processedImg)
	default:
		err = png.Encode(&buf, processedImg)
	}
	if err != nil {
		return nil, fmt.Errorf("文件编码失败: %w", err)
	}
	return buf.Bytes(), nil
}

func (p *Processor) hasTransparency(file image.Image) bool {
	bounds := file.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := file.At(x, y).RGBA()
			if a < 0xffff {
				return true
			}
		}
	}
	return false
}

func (p *Processor) applyImageWatermark(dst *image.RGBA, config *WatermarkConfig) error {
	watermarkImg, err := p.loadWatermarkImage(config)
	if err != nil {
		return fmt.Errorf("加载水印文件失败: %w", err)
	}

	wmBounds := watermarkImg.Bounds()
	originalWidth := wmBounds.Dx()
	originalHeight := wmBounds.Dy()

	// 前端已经生成了缩放后的水印图片，后端不再缩放，直接使用原始尺寸
	// 如果 Scale 未设置或为 0，使用 1.0（不缩放）
	scale := config.Scale
	if scale <= 0 {
		scale = 1.0
	}
	wmWidth := int(float64(originalWidth) * scale)
	wmHeight := int(float64(originalHeight) * scale)
	if wmWidth <= 0 {
		wmWidth = originalWidth
	}
	if wmHeight <= 0 {
		wmHeight = originalHeight
	}

	scaled := p.scaleImageNearest(watermarkImg, wmWidth, wmHeight)
	finalSrc := image.Image(scaled)

	var drawRect image.Rectangle
	if config.Rotation != 0 {
		rotated := p.rotateImageGeneric(scaled, float64(config.Rotation))
		finalSrc = rotated
		rotW := rotated.Bounds().Dx()
		rotH := rotated.Bounds().Dy()
		bounds := dst.Bounds()
		// 使用旋转后的尺寸计算位置，避免被裁切；新算法支持锚点+偏移
		pos := p.calculatePositionWithConfig(bounds.Max.X, bounds.Max.Y, rotW, rotH, config)
		drawRect = image.Rect(pos.X, pos.Y, pos.X+rotW, pos.Y+rotH)
	} else {
		bounds := dst.Bounds()
		pos := p.calculatePositionWithConfig(bounds.Max.X, bounds.Max.Y, wmWidth, wmHeight, config)
		drawRect = image.Rect(pos.X, pos.Y, pos.X+wmWidth, pos.Y+wmHeight)
	}

	if config.Shadow {
		shadowRect := image.Rect(drawRect.Min.X+config.ShadowOffsetX, drawRect.Min.Y+config.ShadowOffsetY, drawRect.Max.X+config.ShadowOffsetX, drawRect.Max.Y+config.ShadowOffsetY)
		// 解析前端传入的阴影颜色，默认黑色
		shadowCol := color.RGBA{0, 0, 0, 255}
		if c, err := p.parseColor(config.ShadowColor); err == nil {
			if cr, ok := c.(color.RGBA); ok {
				shadowCol = cr
			} else {
				r, g, b, a := c.RGBA()
				shadowCol = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			}
		}
		// 提高阴影可见度：与水印不透明度相关联，取其 0.5 倍
		shadowOpacity := config.Opacity * 0.5
		if shadowOpacity < 0.15 {
			shadowOpacity = 0.15
		}
		if shadowOpacity > 1.0 {
			shadowOpacity = 1.0
		}
		p.drawShadowWithMask(dst, finalSrc, shadowRect, shadowCol, shadowOpacity)
	}

	p.drawImageWithOpacity(dst, finalSrc, drawRect, config.Opacity)
	return nil
}

func (p *Processor) drawShadowWithMask(dst *image.RGBA, src image.Image, rect image.Rectangle, shadow color.RGBA, opacity float64) {
	if opacity <= 0 {
		return
	}
	srcBounds := src.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if x < 0 || y < 0 || x >= dst.Bounds().Max.X || y >= dst.Bounds().Max.Y {
				continue
			}
			scaleX := float64(srcBounds.Dx()) / float64(rect.Dx())
			scaleY := float64(srcBounds.Dy()) / float64(rect.Dy())
			srcX := srcBounds.Min.X + int(float64(x-rect.Min.X)*scaleX)
			srcY := srcBounds.Min.Y + int(float64(y-rect.Min.Y)*scaleY)
			if srcX >= srcBounds.Max.X || srcY >= srcBounds.Max.Y {
				continue
			}
			var srcA uint8
			if nrgbaImg, ok := src.(*image.NRGBA); ok {
				srcA = nrgbaImg.NRGBAAt(srcX, srcY).A
			} else {
				_, _, _, a := src.At(srcX, srcY).RGBA()
				srcA = uint8(a >> 8)
			}
			if srcA == 0 {
				continue
			}
			finalAlpha := float64(srcA) / 255.0 * opacity
			if finalAlpha <= 0 {
				continue
			}
			invAlpha := 1.0 - finalAlpha
			dstColor := dst.RGBAAt(x, y)
			newR := uint8(float64(shadow.R)*finalAlpha + float64(dstColor.R)*invAlpha)
			newG := uint8(float64(shadow.G)*finalAlpha + float64(dstColor.G)*invAlpha)
			newB := uint8(float64(shadow.B)*finalAlpha + float64(dstColor.B)*invAlpha)
			newA := uint8(math.Max(finalAlpha*255, float64(dstColor.A)))
			dst.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
}

func (p *Processor) drawImageWithOpacity(dst *image.RGBA, src image.Image, rect image.Rectangle, opacity float64) {
	srcBounds := src.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if x >= 0 && x < dst.Bounds().Max.X && y >= 0 && y < dst.Bounds().Max.Y {
				scaleX := float64(srcBounds.Dx()) / float64(rect.Dx())
				scaleY := float64(srcBounds.Dy()) / float64(rect.Dy())
				srcX := srcBounds.Min.X + int(float64(x-rect.Min.X)*scaleX)
				srcY := srcBounds.Min.Y + int(float64(y-rect.Min.Y)*scaleY)
				if srcX < srcBounds.Max.X && srcY < srcBounds.Max.Y {
					var srcR, srcG, srcB, srcA uint8
					if nrgbaImg, ok := src.(*image.NRGBA); ok {
						nrgbaColor := nrgbaImg.NRGBAAt(srcX, srcY)
						srcR, srcG, srcB, srcA = nrgbaColor.R, nrgbaColor.G, nrgbaColor.B, nrgbaColor.A
					} else {
						r, g, b, a := src.At(srcX, srcY).RGBA()
						srcR, srcG, srcB, srcA = uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
					}
					if srcA == 0 {
						continue
					}
					dstColor := dst.RGBAAt(x, y)
					finalAlpha := float64(srcA) / 255.0 * opacity
					invAlpha := 1.0 - finalAlpha
					newR := uint8(float64(srcR)*finalAlpha + float64(dstColor.R)*invAlpha)
					newG := uint8(float64(srcG)*finalAlpha + float64(dstColor.G)*invAlpha)
					newB := uint8(float64(srcB)*finalAlpha + float64(dstColor.B)*invAlpha)
					newA := uint8(math.Max(finalAlpha*255, float64(dstColor.A)))
					dst.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: newA})
				}
			}
		}
	}
}

func (p *Processor) loadWatermarkImage(config *WatermarkConfig) (image.Image, error) {
	if config.FileBase64 != "" {
		return p.loadFromBase64(config.FileBase64)
	}
	if config.GeneratedImage != "" {
		return p.loadFromBase64(config.GeneratedImage)
	}
	if config.GeneratedFile != "" {
		return p.loadFromBase64(config.GeneratedFile)
	}
	if len(config.FileData) > 0 {
		file, _, err := image.Decode(bytes.NewReader(config.FileData))
		return file, err
	}
	if config.FileURL != "" {
		return p.loadFromFile(config.FileURL)
	}
	return nil, fmt.Errorf("无法加载水印文件：未提供任何水印数据源")
}

func (p *Processor) loadFromBase64(base64Data string) (image.Image, error) {
	if strings.HasPrefix(base64Data, "data:") {
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		} else {
			return nil, fmt.Errorf("无效的base64 data URL格式")
		}
	}
	base64Data = strings.ReplaceAll(base64Data, " ", "")
	base64Data = strings.ReplaceAll(base64Data, "\n", "")
	base64Data = strings.ReplaceAll(base64Data, "\r", "")
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	file, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("文件解码失败: %w", err)
	}
	return file, nil
}

func (p *Processor) loadFromFile(fileURL string) (image.Image, error) {
	filePath := strings.TrimPrefix(fileURL, "/")

	// 安全检查：防止路径遍历攻击
	if strings.Contains(filePath, "..") {
		return nil, fmt.Errorf("非法文件路径：包含 '..' 路径遍历")
	}

	fullPath := fmt.Sprintf("%s/%s", p.watermarkImagePath, filePath)
	cleanPath := strings.ReplaceAll(fullPath, "\\", "/")

	absBasePath := p.watermarkImagePath
	if !strings.HasPrefix(cleanPath, absBasePath) {
		return nil, fmt.Errorf("非法文件路径：超出允许目录范围")
	}

	fileData, err := readFileBytes(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("读取水印文件失败 %s: %w", cleanPath, err)
	}

	img, _, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("解码水印图片失败 %s: %w", cleanPath, err)
	}

	return img, nil
}

func readFileBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	return data, nil
}

// calculatePositionWithConfig 统一的位置计算函数
// 语义：offsetX/offsetY 始终表示"距离参考边缘的距离"（类似CSS margin）
func (p *Processor) calculatePositionWithConfig(imgW, imgH, wmW, wmH int, cfg *WatermarkConfig) Position {
	offsetX := cfg.OffsetX
	offsetY := cfg.OffsetY
	if cfg.OffsetUnit == "percent" {
		offsetX = offsetX * float64(imgW)
		offsetY = offsetY * float64(imgH)
	}

	var x, y int

	switch cfg.Position {
	case PositionTopLeft:
		x = int(offsetX)
		y = int(offsetY)

	case PositionTopCenter:
		// 水平居中，距离顶部
		x = (imgW-wmW)/2 + int(offsetX)
		y = int(offsetY)

	case PositionTopRight:
		x = imgW - wmW - int(offsetX)
		y = int(offsetY)

	case PositionMiddleLeft:
		// 垂直居中，距离左侧
		x = int(offsetX)
		y = (imgH-wmH)/2 + int(offsetY)

	case PositionMiddleCenter:
		// 完全居中，偏移量作为微调
		x = (imgW-wmW)/2 + int(offsetX)
		y = (imgH-wmH)/2 + int(offsetY)

	case PositionMiddleRight:
		// 垂直居中，距离右侧
		x = imgW - wmW - int(offsetX)
		y = (imgH-wmH)/2 + int(offsetY)

	case PositionBottomLeft:
		x = int(offsetX)
		y = imgH - wmH - int(offsetY)

	case PositionBottomCenter:
		// 水平居中，距离底部
		x = (imgW-wmW)/2 + int(offsetX)
		y = imgH - wmH - int(offsetY)

	case PositionBottomRight:
		x = imgW - wmW - int(offsetX)
		y = imgH - wmH - int(offsetY)

	case PositionCustom:
		// 自定义位置（保留，但不推荐使用）
		// 当用户需要完全自定义时，可以通过前端直接指定具体坐标
		x = int(offsetX)
		y = int(offsetY)

	default:
		x = imgW - wmW - int(offsetX)
		y = imgH - wmH - int(offsetY)
	}

	// 边界限制（允许负值，实现水印部分超出图片）
	// 但确保至少有一部分在图片内
	if x < -wmW {
		x = -wmW + 1
	}
	if y < -wmH {
		y = -wmH + 1
	}
	if x > imgW {
		x = imgW - 1
	}
	if y > imgH {
		y = imgH - 1
	}

	return Position{X: x, Y: y}
}

func (p *Processor) parseColor(colorStr string) (color.Color, error) {
	if colorStr == "" {
		return nil, fmt.Errorf("颜色字符串为空")
	}
	colorStr = strings.TrimPrefix(colorStr, "#")
	if len(colorStr) == 3 {
		colorStr = string(colorStr[0]) + string(colorStr[0]) +
			string(colorStr[1]) + string(colorStr[1]) +
			string(colorStr[2]) + string(colorStr[2])
	}
	if len(colorStr) != 6 {
		return nil, fmt.Errorf("无效的颜色格式: %s", colorStr)
	}
	rgb, err := hex.DecodeString(colorStr)
	if err != nil {
		return nil, fmt.Errorf("颜色解析失败: %w", err)
	}
	return color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255}, nil
}

func (p *Processor) ValidateConfig(config *WatermarkConfig) error {
	if config == nil {
		return fmt.Errorf("配置不能为空")
	}
	if !config.Enabled {
		return nil
	}
	if config.Type != TypeImage && config.Type != TypeFile {
		return fmt.Errorf("仅支持图片水印")
	}

	if config.FileBase64 == "" && len(config.FileData) == 0 && config.GeneratedImage == "" && config.GeneratedFile == "" && config.FileURL == "" {
		return fmt.Errorf("图片水印必须指定水印数据源（fileBase64[前端生成]/fileURL[后端文件]/其他）")
	}
	if config.Opacity < 0 || config.Opacity > 1 {
		return fmt.Errorf("透明度必须在0-1之间")
	}
	if config.OffsetUnit != "" && config.OffsetUnit != "px" && config.OffsetUnit != "percent" {
		return fmt.Errorf("OffsetUnit 仅支持 px 或 percent")
	}
	if config.OffsetUnit == "percent" {
		if config.OffsetX < 0 || config.OffsetX > 1 || config.OffsetY < 0 || config.OffsetY > 1 {
			return fmt.Errorf("当 OffsetUnit=percent 时，offsetX/offsetY 必须在 0-1 之间")
		}
	}
	return nil
}

func (p *Processor) GetSupportedFormats() []string { return []string{"jpeg", "jpg", "png"} }

func CreateWatermarkResult(success bool, file image.Image, data []byte, err error, processTime int64) *WatermarkResult {
	result := &WatermarkResult{
		Success:        success,
		ProcessedImage: file,
		ProcessedData:  data,
		ProcessTime:    processTime,
	}

	if err != nil {
		result.ErrorMessage = err.Error()
	}

	return result
}

func (p *Processor) scaleImageNearest(src image.Image, dstW, dstH int) *image.RGBA {
	srcB := src.Bounds()
	srcW := srcB.Dx()
	srcH := srcB.Dy()
	if srcW == 0 || srcH == 0 || dstW <= 0 || dstH <= 0 {
		return image.NewRGBA(image.Rect(0, 0, maxInt(1, dstW), maxInt(1, dstH)))
	}
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	scaleX := float64(srcW) / float64(dstW)
	scaleY := float64(srcH) / float64(dstH)
	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			sx := srcB.Min.X + int(float64(x)*scaleX)
			sy := srcB.Min.Y + int(float64(y)*scaleY)
			c := src.At(sx, sy)
			dst.Set(x, y, c)
		}
	}
	return dst
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (p *Processor) rotateImageGeneric(src image.Image, degrees float64) *image.RGBA {
	if degrees == 0 {
		bounds := src.Bounds()
		dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Src)
		return dst
	}
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	newSize := int(math.Sqrt(float64(width*width+height*height))) + 10
	dst := image.NewRGBA(image.Rect(0, 0, newSize, newSize))
	centerX := float64(newSize / 2)
	centerY := float64(newSize / 2)
	srcCenterX := float64(width / 2)
	srcCenterY := float64(height / 2)
	rad := degrees * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	for y := 0; y < newSize; y++ {
		for x := 0; x < newSize; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			srcX := dx*cos + dy*sin + srcCenterX
			srcY := -dx*sin + dy*cos + srcCenterY
			if srcX >= 0 && srcX < float64(width) && srcY >= 0 && srcY < float64(height) {
				c := src.At(int(srcX)+bounds.Min.X, int(srcY)+bounds.Min.Y)
				dst.Set(x, y, c)
			}
		}
	}
	return dst
}
