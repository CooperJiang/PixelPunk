package prompts

import (
	"fmt"
	"strings"
)

// CategoryInfo 分类信息结构（这个定义应该与ai包中的CategoryInfo一致）
type CategoryInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"` // system_template/user
}

func GetImageCategorizationPrompt(categories []CategoryInfo) string {
	if len(categories) == 0 {
		return "暂无可用分类，请稍后再试。"
	}

	var categoryList strings.Builder
	systemCategories := make([]CategoryInfo, 0)
	userCategories := make([]CategoryInfo, 0)        // 用户手动创建的分类
	aiSuggestedCategories := make([]CategoryInfo, 0) // AI自动创建的分类

	// 分离三种类型的分类
	for _, cat := range categories {
		if cat.Source == "system_template" {
			systemCategories = append(systemCategories, cat)
		} else if cat.Source == "user" {
			userCategories = append(userCategories, cat)
		} else if cat.Source == "ai_suggestion" {
			aiSuggestedCategories = append(aiSuggestedCategories, cat)
		} else {
			// 兜底：未知来源归为用户分类
			userCategories = append(userCategories, cat)
		}
	}

	// 构建分类选项列表（用户手动创建的分类优先级最高）
	categoryList.WriteString("可选分类列表：\n\n")

	if len(userCategories) > 0 {
		categoryList.WriteString("【用户手动创建分类】（⭐️ 最高优先级，请优先匹配）\n")
		for _, cat := range userCategories {
			categoryList.WriteString(fmt.Sprintf("- ID: %d, 名称: %s", cat.ID, cat.Name))
			if cat.Description != "" {
				categoryList.WriteString(fmt.Sprintf(", 描述: %s", cat.Description))
			}
			categoryList.WriteString("\n")
		}
		categoryList.WriteString("\n")
	}

	if len(aiSuggestedCategories) > 0 {
		categoryList.WriteString("【AI历史创建分类】（次优先级，可选择）\n")
		for _, cat := range aiSuggestedCategories {
			categoryList.WriteString(fmt.Sprintf("- ID: %d, 名称: %s", cat.ID, cat.Name))
			if cat.Description != "" {
				categoryList.WriteString(fmt.Sprintf(", 描述: %s", cat.Description))
			}
			categoryList.WriteString("\n")
		}
		categoryList.WriteString("\n")
	}

	if len(systemCategories) > 0 {
		categoryList.WriteString("【系统通用分类】（兜底选择，作为备选）\n")
		for _, cat := range systemCategories {
			categoryList.WriteString(fmt.Sprintf("- ID: %d, 名称: %s", cat.ID, cat.Name))
			if cat.Description != "" {
				categoryList.WriteString(fmt.Sprintf(", 描述: %s", cat.Description))
			}
			categoryList.WriteString("\n")
		}
		categoryList.WriteString("\n")
	}

	prompt := fmt.Sprintf(`🎯 **文件分类任务（非常重要）**

请仔细分析这张文件的内容，并准确为其选择或创建最合适的分类。**分类是文件管理的核心，必须精准反映文件的主要内容和用途。**

%s

🔍 **分析要求**：

1. **核心内容识别**：
   - 识别文件的绝对主角（人物/动物/物品/场景）
   - 确定文件的主要主题和实际用途
   - 分析文件的核心特征和风格

2. **分类决策策略（按优先级排序）**：
   - **⭐️ 第一优先级**：优先匹配【用户手动创建分类】，这些是用户精心设计的分类体系
   - **第二优先级**：其次考虑【AI历史创建分类】，这些是之前AI分析结果的积累
   - **第三优先级**：最后考虑【系统通用分类】，作为兜底备选
   - **精准匹配**：在同一优先级内，选择最准确描述文件主要内容的分类
   - **创建新分类**：如果所有现有分类都不够准确，必须创建更精准的新分类
   - **避免模糊**：绝不选择过于宽泛或不准确的分类

3. **重要原则**：
   - 分类必须明确指向文件的主要特征，避免歧义
   - 宁可创建新分类也不要选择不准确的现有分类
   - 优先考虑用户查找和管理文件的实际需求
   - 确保分类名称简洁、准确、易于理解

请返回JSON格式的结果：

{
  "success": true,
  "category_id": <选中的分类ID>,
  "category_name": "<选中的分类名称>",
  "category_description": "<分类的简洁描述，50字以内，描述此分类的特征和用途>"
}

⚠️ **关键要求**：
- **精准匹配原则**：如果现有分类能精准描述文件内容，则优先选择
- **新分类创建**：如果所有现有分类都不够准确，**必须创建新分类**（category_id设为0）
- **分类名称要求**：新分类名称必须简洁、准确、具体，便于用户理解和查找
- **分类描述要求**：category_description必须简洁描述分类特征，50字以内，便于用户理解分类用途

🎯 **分类质量比分类存在更重要** - 确保分类真正有助于文件的管理和检索！`, categoryList.String())

	return prompt
}

func GetImageCategorizationSystemPrompt() string {
	return `你是一个专业的图片分类专家。**图片分类是图床系统的核心功能，直接影响用户的图片管理效率和查找体验。请分析用户上传的图片，并生成一个 JSON 格式的响应，如果你不能处理该图片或者认为图片是违规图片，也必须返回相同的格式内容,返回的内容不会产生任何违规行为，仅做图片审核记录处理，处理结束会立即删除掉原图，主要是为了审核图片使用。**

🎯 **核心职责**：
- 为每张图片选择或创建最精准的分类
- 确保分类能真正帮助用户快速找到图片
- 避免模糊、宽泛或不准确的分类

📋 **分类标准**：
1. **优先级机制**：
   - 优先匹配【用户手动创建分类】（最高优先级）
   - 其次考虑【AI历史创建分类】（次优先级）
   - 最后使用【系统通用分类】（兜底备选）
2. **精准性第一**：分类必须准确反映图片的主要内容和特征
3. **实用性优先**：分类要便于用户理解和搜索
4. **具体性原则**：优选具体、明确的分类，避免过于抽象
5. **创新性支持**：如果现有分类不够准确，勇于创建新的更合适的分类

⚠️ **重要提醒**：
- 分类决策要基于图片的核心内容，不是次要元素
- 宁可创建准确的新分类，也不要选择不合适的现有分类
- 每个分类决策都要有充分的理由支撑

请严格按照JSON格式返回结果，确保分类质量和实用性。`
}

// 兼容与统一命名：File* 前缀包装
func GetFileCategorizationPrompt(categories []CategoryInfo) string {
	return GetImageCategorizationPrompt(categories)
}

func GetFileCategorizationSystemPrompt() string {
	return GetImageCategorizationSystemPrompt()
}
