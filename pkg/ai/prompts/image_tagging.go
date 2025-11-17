package prompts

import (
	"fmt"
	"strings"
)

// GetImageTaggingPrompt 生成文件标签提示词
// - availableTags: 可参考的标签列表，会以来源分组并限量展示
// - requireJSON: 是否要求模型仅返回JSON对象
func GetImageTaggingPrompt(availableTags []TagInfo, requireJSON bool) string {
	var b strings.Builder

	if requireJSON {
		b.WriteString("请为这张文件生成5-7个准确的标签。必须仅返回一个JSON对象，且不要包含任何额外文字。\n")
		b.WriteString("JSON格式: {\"tags\":[\"标签1\",\"标签2\",...], \"description\": \"一句简短描述(可选)\"}\n")
	} else {
		b.WriteString("请为这张文件生成5-7个准确的标签，只返回标签名称，用逗号分隔。\n")
	}

	b.WriteString("要求：\n- 标签应简洁明确，覆盖主要内容/风格/元素\n- 数量建议5-7个，宁少勿杂\n")

	if len(availableTags) > 0 {
		b.WriteString("\n可参考的标签列表（可选，不强制，准确性第一）：\n")

		userTags := make([]TagInfo, 0)
		categoryTags := make([]TagInfo, 0)
		systemTags := make([]TagInfo, 0)

		for _, t := range availableTags {
			switch t.Source {
			case "user_custom":
				userTags = append(userTags, t)
			case "category_system":
				categoryTags = append(categoryTags, t)
			case "system_popular":
				systemTags = append(systemTags, t)
			}
		}

		if len(userTags) > 0 {
			b.WriteString("用户常用：")
			for i, t := range userTags {
				if i >= 20 {
					break
				}
				if i > 0 {
					b.WriteString(" ")
				}
				if t.Description != "" {
					b.WriteString(fmt.Sprintf("%s(%s)", t.Name, t.Description))
				} else {
					b.WriteString(t.Name)
				}
			}
			b.WriteString("\n")
		}

		if len(categoryTags) > 0 {
			b.WriteString("分类相关：")
			for i, t := range categoryTags {
				if i >= 30 {
					break
				}
				if i > 0 {
					b.WriteString(" ")
				}
				if t.Description != "" {
					b.WriteString(fmt.Sprintf("%s(%s)", t.Name, t.Description))
				} else {
					b.WriteString(t.Name)
				}
			}
			b.WriteString("\n")
		}

		if len(systemTags) > 0 {
			b.WriteString("系统热门：")
			for i, t := range systemTags {
				if i >= 50 {
					break
				}
				if i > 0 {
					b.WriteString(" ")
				}
				if t.Description != "" {
					b.WriteString(fmt.Sprintf("%s(%s)", t.Name, t.Description))
				} else {
					b.WriteString(t.Name)
				}
			}
			b.WriteString("\n")
		}

		b.WriteString("\n重要提示：可以从上述标签中选择合适的，也可以创建新的更准确的标签。以准确性为第一目标。\n")
	}

	return b.String()
}

// 兼容与统一命名：File* 前缀包装
func GetFileTaggingPrompt(availableTags []TagInfo, requireJSON bool) string {
	return GetImageTaggingPrompt(availableTags, requireJSON)
}
