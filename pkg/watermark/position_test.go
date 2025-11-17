package watermark

import (
	"testing"
)

// TestPositionCalculation 验证9个锚点位置的计算逻辑
func TestPositionCalculation(t *testing.T) {
	processor := NewProcessor()

	imgW, imgH := 1000, 800
	wmW, wmH := 100, 50
	offsetX, offsetY := 20.0, 20.0

	tests := []struct {
		name        string
		position    WatermarkPosition
		expectedX   int
		expectedY   int
		description string
	}{
		{
			name:        "top-left",
			position:    PositionTopLeft,
			expectedX:   20, // offsetX
			expectedY:   20, // offsetY
			description: "距离左上角各20px",
		},
		{
			name:        "top-center",
			position:    PositionTopCenter,
			expectedX:   (imgW-wmW)/2 + 20, // 居中 + 偏移
			expectedY:   20,
			description: "水平居中，距离顶部20px",
		},
		{
			name:        "top-right",
			position:    PositionTopRight,
			expectedX:   imgW - wmW - 20, // 距离右边缘20px
			expectedY:   20,
			description: "距离右上角各20px",
		},
		{
			name:        "middle-left",
			position:    PositionMiddleLeft,
			expectedX:   20,
			expectedY:   (imgH-wmH)/2 + 20, // 垂直居中 + 偏移
			description: "垂直居中，距离左侧20px",
		},
		{
			name:        "middle-center",
			position:    PositionMiddleCenter,
			expectedX:   (imgW-wmW)/2 + 20,
			expectedY:   (imgH-wmH)/2 + 20,
			description: "完全居中，偏移各20px作为微调",
		},
		{
			name:        "middle-right",
			position:    PositionMiddleRight,
			expectedX:   imgW - wmW - 20,
			expectedY:   (imgH-wmH)/2 + 20,
			description: "垂直居中，距离右侧20px",
		},
		{
			name:        "bottom-left",
			position:    PositionBottomLeft,
			expectedX:   20,
			expectedY:   imgH - wmH - 20, // 距离底边缘20px
			description: "距离左下角各20px",
		},
		{
			name:        "bottom-center",
			position:    PositionBottomCenter,
			expectedX:   (imgW-wmW)/2 + 20,
			expectedY:   imgH - wmH - 20,
			description: "水平居中，距离底部20px",
		},
		{
			name:        "bottom-right",
			position:    PositionBottomRight,
			expectedX:   imgW - wmW - 20, // 距离右边缘20px
			expectedY:   imgH - wmH - 20, // 距离底边缘20px
			description: "距离右下角各20px",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &WatermarkConfig{
				Position:   tt.position,
				OffsetX:    offsetX,
				OffsetY:    offsetY,
				OffsetUnit: "px",
			}

			pos := processor.calculatePositionWithConfig(imgW, imgH, wmW, wmH, cfg)

			if pos.X != tt.expectedX {
				t.Errorf("%s: X坐标错误\n  期望: %d (距离参考边缘%dpx)\n  实际: %d\n  说明: %s",
					tt.name, tt.expectedX, int(offsetX), pos.X, tt.description)
			}

			if pos.Y != tt.expectedY {
				t.Errorf("%s: Y坐标错误\n  期望: %d (距离参考边缘%dpx)\n  实际: %d\n  说明: %s",
					tt.name, tt.expectedY, int(offsetY), pos.Y, tt.description)
			}

			t.Logf("✓ %s 通过: (%d, %d) - %s", tt.name, pos.X, pos.Y, tt.description)
		})
	}
}

// TestPercentOffset 验证百分比偏移
func TestPercentOffset(t *testing.T) {
	processor := NewProcessor()

	imgW, imgH := 1000, 800
	wmW, wmH := 100, 50

	cfg := &WatermarkConfig{
		Position:   PositionTopLeft,
		OffsetX:    0.05, // 5%
		OffsetY:    0.05, // 5%
		OffsetUnit: "percent",
	}

	pos := processor.calculatePositionWithConfig(imgW, imgH, wmW, wmH, cfg)

	expectedX := int(0.05 * float64(imgW)) // 50px
	expectedY := int(0.05 * float64(imgH)) // 40px

	if pos.X != expectedX || pos.Y != expectedY {
		t.Errorf("百分比偏移计算错误: 期望(%d, %d), 实际(%d, %d)",
			expectedX, expectedY, pos.X, pos.Y)
	} else {
		t.Logf("✓ 百分比偏移通过: 5%% = (%dpx, %dpx)", pos.X, pos.Y)
	}
}

// TestNegativeOffset 验证负偏移（水印部分超出边界）
func TestNegativeOffset(t *testing.T) {
	processor := NewProcessor()

	imgW, imgH := 1000, 800
	wmW, wmH := 100, 50

	cfg := &WatermarkConfig{
		Position:   PositionTopLeft,
		OffsetX:    -10, // 向左伸出10px
		OffsetY:    -10, // 向上伸出10px
		OffsetUnit: "px",
	}

	pos := processor.calculatePositionWithConfig(imgW, imgH, wmW, wmH, cfg)

	// 允许负值，但不能完全超出
	if pos.X >= 0 || pos.Y >= 0 {
		t.Errorf("负偏移应该允许部分超出边界: 实际(%d, %d)", pos.X, pos.Y)
	}

	t.Logf("✓ 负偏移允许部分超出: (%d, %d)", pos.X, pos.Y)
}
