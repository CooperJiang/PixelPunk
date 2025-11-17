// Package watermark provides file watermarking functionality
// 水印处理包 - 提供文件水印处理功能
package watermark

import (
	"image"
)

// 全局变量 - 默认水印管理器
var defaultManager *Manager

// 包初始化
func init() {
	defaultManager = NewManager(DefaultManagerConfig())
}

func SetDefaultManager(manager *Manager) {
	defaultManager = manager
}

func GetDefaultManager() *Manager {
	return defaultManager
}

// 便捷函数 - 使用默认管理器处理

// ProcessImage 使用默认管理器处理图像
func ProcessImage(img image.Image, config *WatermarkConfig) (*WatermarkResult, error) {
	return defaultManager.Process(img, config)
}

// ProcessImageBytes 使用默认管理器处理图像字节数据
func ProcessImageBytes(imgData []byte, config *WatermarkConfig) (*WatermarkResult, error) {
	return defaultManager.ProcessBytes(imgData, config)
}

// ProcessWithConfigJSON 使用JSON配置处理图像
func ProcessWithConfigJSON(img image.Image, configJSON string) (*WatermarkResult, error) {
	return defaultManager.ProcessWithConfig(img, configJSON)
}

// ProcessBytesWithConfigJSON 使用JSON配置处理图像字节数据
func ProcessBytesWithConfigJSON(imgData []byte, configJSON string) (*WatermarkResult, error) {
	return defaultManager.ProcessBytesWithConfig(imgData, configJSON)
}

// ValidateConfig 验证水印配置
func ValidateConfig(config *WatermarkConfig) error {
	return defaultManager.ValidateConfig(config)
}

// ValidateConfigJSON 验证JSON配置字符串
func ValidateConfigJSON(configJSON string) error {
	return defaultManager.ValidateConfigJSON(configJSON)
}

// IsConfigValid 检查配置是否有效
func IsConfigValid(configJSON string) bool {
	return defaultManager.IsConfigValid(configJSON)
}

func GetStats() *ProcessingStats {
	return defaultManager.GetStats()
}

// ResetStats 重置统计信息
func ResetStats() {
	defaultManager.ResetStats()
}

func GetSupportedFormats() []string {
	return defaultManager.GetSupportedFormats()
}

func SetWatermarkPath(path string) {
	config := defaultManager.GetConfig()
	config.WatermarkPath = path
	defaultManager.SetConfig(config)
}

// Example 使用示例
/*
package main

import (
	"fmt"
	"pixelpunk/pkg/watermark"
)

func main() {
	config := &watermark.WatermarkConfig{
		Enabled:   true,
		Type:      watermark.TypeText,
		Text:      "© PixelPunk 2024",
		FontSize:  24,
		FontColor: "#ffffff",
		Position:  watermark.PositionBottomRight,
		Opacity:   0.7,
		Shadow:    true,
	}

	// 处理图像（假设已有图像数据）
	result, err := watermark.ProcessImageBytes(imageData, config)
	if err != nil {
		return
	}

	if result.Success {
		// 使用 result.ProcessedData 作为处理后的图像数据
	}
}
*/
