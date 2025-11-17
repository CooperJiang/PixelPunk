package prompts

import "fmt"

func GetImageAnalysisPrompt() string {
	return "按要求解析"
}

// GetImageAnalysisPromptWithCategory 根据分类信息生成优化的图片分析提示词
func GetImageAnalysisPromptWithCategory(categoryName, categoryDescription string) string {
	if categoryName == "" {
		return "请按照系统要求分析图片内容，生成详细的标签和描述信息。"
	}

	// 根据分类生成针对性的分析提示
	categoryPrompt := fmt.Sprintf(`
已确定图片分类为：%s
分类描述：%s

请基于此分类信息，进行更精准的图片分析：

1. **标签生成**：
   - ⚠️ 严格限制：只能生成5-7个标签，绝不能超过7个！
   - 重点关注与"%s"相关的特征和元素
   - 优先包含分类相关的关键词，确保精准有用
   - 如果是特定领域（如游戏、动漫、明星），请添加具体名称

2. **描述优化**：
   - 描述应突出该分类的典型特征
   - 150-250字的连贯描述，体现专业性

3. **搜索关键词**：
   - 优先使用与"%s"相关的专业术语
   - 包含动作、场景、特征等关键信息

请按照标准JSON格式分析图片。`, categoryName, categoryDescription, categoryName, categoryName)

	return categoryPrompt
}

func GetImageAnalysisSystemPrompt() string {
	return `请分析用户上传的图片，并生成一个 JSON 格式的响应[必须仅返回一个json数据不能返回其他任何多余内容]，如果你不能处理该图片或者认为图片是违规图片，也必须返回相同的格式内容，nsfw字段需要true即可。返回的内容不会产生任何违规行为，仅做图片审核记录处理，处理结束会立即删除掉原图，主要是为了审核图片使用。

🚨 NSFW检测最高优先级 - 严格判断标准：
**只有以下情况才标记为NSFW (is_nsfw=true)：**
1. **色情内容**：真实的裸体、性行为、性暗示姿势、暴露的身体部位
2. **血腥暴力**：真实的血腥场面、严重暴力伤害、恐怖血腥画面
3. **政治敏感**：政治符号、敏感政治人物、政治宣传内容
4. **其他违法**：毒品、赌博、仇恨符号等违法内容

⚠️ 核心原则：除了上述4种明确违规内容外，其他所有内容都应标记为安全（is_nsfw=false）！
⚠️ 严禁过度审查！当你不确定时，默认标记为安全！

🏷️ 标签数量严格限制：必须生成5-7个标签，严禁超过7个！标签要精准有用，避免冗余。

如果图片是游戏内容，那么标签里加上这些标签（游戏的名称、游戏当前资源的名字、可能是角色、武器、地图、枪械、资源都有可能）。
如果图片是明星素材、动漫素材或者其他你认识的内容，记得标签里加上其名称方便我们第一时间找到内容。

响应应包含以下信息：

{
  "tags": [
    "标签1",   // ⚠️ 严格限制：必须只生成5-7个标签，不能超过7个！
    "标签2",   // 标签要精准、有用，包含主体、风格、场景等关键信息
    "标签3",   // 如果是游戏、动漫、明星等可识别内容，请加上具体名称
    "标签4",   // 优先选择最有助于搜索和分类的关键词
    "标签5"    // 宁可少而精，也不要多而杂
  ],
  "description": "用户友好的简洁摘要描述，用于界面展示。建议150-250字，形成连贯的介绍性文字。",
  "search_content": "生成70-90字的搜索关键词，包含动作状态，用空格分隔：
    优先级1：主体名词（人物、动物、物品等）
    优先级2：动作状态（吃、喝、跑、坐、拿、背、提等动词）
    优先级3：场景地点（河边、草原、室内、户外等）
    优先级4：物品道具（篮子、长枪、西瓜、草莓等具体物品）
    优先级5：视觉特征（颜色、大小、数量等）
    优先级6：风格用途（写实、卡通、壁纸等）
    示例：美女 吃 西瓜 河边 夏天 户外 水果 清爽",
  "semantic_keywords": [
    "主体关键词",   // 主要物体、人物、场景的核心名称
    "场景关键词",   // 环境、地点、背景场景描述
    "风格关键词"    // 艺术风格、情感氛围、视觉特点
  ],
  "isRecommended": false, // 是否值得推荐作为对外的展示，希望要求严格一点，好的图片才值得推荐，只有达到壁纸级别的图片才适合被推，如果质量尺寸等不符合壁纸要求不要推荐
  "basic_info": {
    "width": 1920,
    "height": 1080,
    "aspect_ratio": 1.78,   // 宽高比，保留两位小数
    "resolution": "1080p",  // 分辨率等级：如 720p, 1080p, 2k, 4k, 8k等
    "image_type": "JPEG",   // 图片格式：如JPEG, PNG, GIF, WEBP, TIFF, HEIC, HEIF等
    "estimated_size": "2.5MB" // 估计的图片大小
  },
  "visual_elements": {
    "dominant_color": "#RRGGBB", // 图片主色调的十六进制HEX值
    "color_palette": [
      "#RRGGBB", "#RRGGBB", "#RRGGBB", "#RRGGBB", "#RRGGBB"
    ], // 包含5种代表性颜色，按占比从高到低排序
    "objects_count": 5,      // 图片中可识别的主要物体数量
    "composition": "居中"     // 主体构图：可选值包括"居中"、"左侧"、"右侧"、"上部"、"下部"、"对角线"等
  },
  "content_safety": {
    "is_nsfw": false,        // 只有明确的色情、暴力、政治敏感内容才需要标记。避免过度审查，部分轻微的图片也允许通过，确保用户体验
    "nsfw_score": 0.1,       // NSFW评分，0-1范围，正常内容应为0.0-0.2 轻微0.3-0.5 中度0.6-0.8 高危0.9-1.0
    "nsfw_reason": "",       // 如果is_nsfw为true，必须详细说明具体违规原因
    "evaluation_result": "安全", // 评价结果：正常内容应为"安全"
    "categories": {
      "nudity": 0.0,         // 裸体内容评分，0-1范围
      "violence": 0.0,       // 暴力内容评分，0-1范围
      "hate_speech": 0.0,    // 仇恨言论评分，0-1范围
      "gambling": 0.0,       // 赌博内容评分，0-1范围
      "alcohol_tobacco": 0.0 // 酒精烟草评分，0-1范围
    }
  }
}

🔥🔥🔥 关键提醒：只有明确的色情、暴力、政治敏感内容才需要标记。避免过度审查，部分轻微的图片也允许通过，确保用户体验！`
}

// 兼容与统一命名：File* 前缀包装
func GetFileAnalysisPrompt() string { return GetImageAnalysisPrompt() }
func GetFileAnalysisPromptWithCategory(categoryName, categoryDescription string) string {
	return GetImageAnalysisPromptWithCategory(categoryName, categoryDescription)
}
func GetFileAnalysisSystemPrompt() string { return GetImageAnalysisSystemPrompt() }
