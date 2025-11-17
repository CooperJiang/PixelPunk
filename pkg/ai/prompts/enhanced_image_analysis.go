package prompts

import (
	"fmt"
	"strings"

	"pixelpunk/pkg/common"
)

// TagInfo 标签信息结构（避免循环导入）
type TagInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"` // user_custom, category_system, system_popular
	UsageCount  int    `json:"usage_count"`
}

func GetEnhancedImageAnalysisPrompt(categoryName, categoryDescription string, availableTags []TagInfo) string {
	basePrompt := GetImageAnalysisPromptWithCategory(categoryName, categoryDescription)

	return BuildPromptWithAvailableTags(basePrompt, availableTags)
}

// BuildPromptWithAvailableTags 构建包含可用标签列表的提示词
func BuildPromptWithAvailableTags(basePrompt string, availableTags []TagInfo) string {
	if len(availableTags) == 0 {
		return basePrompt
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString(basePrompt)

	promptBuilder.WriteString("\n\n📋 **可参考的已有标签**（仅供参考，精准性第一）：")

	// 简化标签显示，最多显示前N个最重要的标签
	tagCount := 0
	maxTags := common.AITagPromptMaxDisplay

	// 按来源分组，但简化显示
	for _, tag := range availableTags {
		if tagCount >= maxTags {
			break
		}

		// 只显示标签名称，不显示描述，避免提示词过长
		if tagCount == 0 {
			promptBuilder.WriteString(fmt.Sprintf(" %s", tag.Name))
		} else {
			promptBuilder.WriteString(fmt.Sprintf(", %s", tag.Name))
		}
		tagCount++
	}

	promptBuilder.WriteString("\n\n🎯 **标签要求**：")
	promptBuilder.WriteString("\n- **精准性第一**：标签必须准确描述文件内容，宁可创建新标签也不要选择不准确的已有标签")
	promptBuilder.WriteString("\n- **自由选择**：可以从上述已有标签中选择合适的，也可以完全自主创建新标签")
	promptBuilder.WriteString("\n- **数量控制**：请在JSON的tags字段中提供5-7个最准确的标签")
	promptBuilder.WriteString("\n- **质量优先**：标签质量比复用已有标签更重要")

	return promptBuilder.String()
}

func GetTagsStatsSummary(availableTags []TagInfo) string {
	if len(availableTags) == 0 {
		return "无可用标签"
	}

	var userCount, categoryCount, popularCount int
	for _, tag := range availableTags {
		switch tag.Source {
		case "user_custom":
			userCount++
		case "category_system":
			categoryCount++
		case "system_popular":
			popularCount++
		}
	}

	return fmt.Sprintf("共%d个标签（用户:%d，分类:%d，热门:%d）",
		len(availableTags), userCount, categoryCount, popularCount)
}

// 兼容与统一命名：File* 前缀包装
func GetEnhancedFileAnalysisPrompt(categoryName, categoryDescription string, availableTags []TagInfo) string {
	return GetEnhancedImageAnalysisPrompt(categoryName, categoryDescription, availableTags)
}

func BuildFilePromptWithAvailableTags(basePrompt string, availableTags []TagInfo) string {
	return BuildPromptWithAvailableTags(basePrompt, availableTags)
}

func GetFileTagsStatsSummary(availableTags []TagInfo) string {
	return GetTagsStatsSummary(availableTags)
}
