package ai

import (
	"encoding/json"
	"fmt"
	"pixelpunk/pkg/ai"
	"strings"
)

// AIFileResponse 是解析后的文件分析结果
type AIFileResponse struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data,omitempty"`
	ErrMsg       string      `json:"errMsg,omitempty"`
	FileURL      string      `json:"fileUrl,omitempty"`
	Usage        *TokenUsage `json:"usage,omitempty"`
	RawResponse  string      `json:"-"`
	HttpDuration int64       `json:"http_duration,omitempty"` // HTTP调用耗时（毫秒）
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func convertToOldResponse(newResp *ai.AIResponse) *AIFileResponse {
	oldResp := &AIFileResponse{
		Success:      newResp.Success,
		ErrMsg:       newResp.ErrMsg,
		FileURL:      newResp.ImageURL,
		Usage:        convertTokenUsage(newResp.Usage),
		RawResponse:  newResp.Data,
		HttpDuration: newResp.HttpDuration,
	}

	if newResp.Success && newResp.Data != "" {
		var dataMap map[string]interface{}
		if err := json.Unmarshal([]byte(newResp.Data), &dataMap); err != nil {
			sample := newResp.Data
			if len(sample) > 500 {
				sample = sample[:500] + "..."
			}
			oldResp.Data = newResp.Data
		} else {
			oldResp.Data = dataMap
		}
	}
	return oldResp
}

// convertTokenUsage 转换Token使用情况
func convertTokenUsage(newUsage *ai.TokenUsage) *TokenUsage {
	if newUsage == nil {
		return nil
	}
	return &TokenUsage{
		PromptTokens:     newUsage.PromptTokens,
		CompletionTokens: newUsage.CompletionTokens,
		TotalTokens:      newUsage.TotalTokens,
	}
}

// parseAITaggingResult 解析AI返回的标记结果为结构化数据
func parseAITaggingResult(data interface{}) (*AITaggingResult, error) {
	var jsonData []byte
	if strData, ok := data.(string); ok {
		// 提取并清洗可能夹杂markdown/思考的JSON
		cleaned := ai.CleanJSON(ai.ExtractJSONFromText(strData))
		cleaned = strings.TrimSpace(cleaned)
		if cleaned == "" || !strings.HasPrefix(cleaned, "{") {
			return nil, fmt.Errorf("AI返回非JSON格式数据")
		}
		jsonData = []byte(cleaned)
	} else {
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("转换AI返回数据为JSON失败: %v", err)
		}
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, fmt.Errorf("解析AI返回数据为Map失败: %v", err)
	}

	if _, has := rawData["content_safety"]; has {
		return parseNewFormatAIResult(rawData, jsonData)
	}
	if _, has := rawData["nsfw_detection"]; has {
		return parseNewFormatAIResult(rawData, jsonData)
	}

	var result AITaggingResult
	if err := json.Unmarshal(jsonData, &result); err != nil {
		// 结构体解析失败，退回到 map 提取关键字段
		result = AITaggingResult{}
	}
	// 补充提取：若结构体字段为空，从 rawData 映射
	if result.Description == "" {
		if s, ok := rawData["description"].(string); ok {
			result.Description = s
		}
	}
	if len(result.Tags) == 0 {
		if arr, ok := rawData["tags"].([]interface{}); ok {
			tags := make([]string, 0, len(arr))
			for _, v := range arr {
				if sv, ok := v.(string); ok {
					tags = append(tags, strings.TrimSpace(sv))
				}
			}
			result.Tags = tags
		}
	}
	if result.SearchContent == "" {
		if s, ok := rawData["search_content"].(string); ok {
			result.SearchContent = s
		}
	}
	return &result, nil
}

// parseNewFormatAIResult 解析新格式的AI标记结果
func parseNewFormatAIResult(rawData map[string]interface{}, jsonData []byte) (*AITaggingResult, error) {
	result := &AITaggingResult{}

	if basicInfo, ok := rawData["basic_info"].(map[string]interface{}); ok {
		if width, ok := basicInfo["width"].(float64); ok {
			result.BasicInfo.Width = int(width)
		}
		if height, ok := basicInfo["height"].(float64); ok {
			result.BasicInfo.Height = int(height)
		}
		if ar, ok := basicInfo["aspect_ratio"].(float64); ok {
			result.BasicInfo.AspectRatio = ar
		} else if result.BasicInfo.Width > 0 && result.BasicInfo.Height > 0 {
			result.BasicInfo.AspectRatio = float64(result.BasicInfo.Width) / float64(result.BasicInfo.Height)
		}
		if s, ok := basicInfo["resolution"].(string); ok {
			result.BasicInfo.Resolution = s
		}
		if s, ok := basicInfo["image_type"].(string); ok {
			result.BasicInfo.ImageType = s
		}
		if s, ok := basicInfo["estimated_size"].(string); ok {
			result.BasicInfo.EstimatedSize = s
		}
	} else {
		if width, ok := rawData["width"].(float64); ok {
			result.BasicInfo.Width = int(width)
		}
		if height, ok := rawData["height"].(float64); ok {
			result.BasicInfo.Height = int(height)
		}
		if result.BasicInfo.Width > 0 && result.BasicInfo.Height > 0 {
			result.BasicInfo.AspectRatio = float64(result.BasicInfo.Width) / float64(result.BasicInfo.Height)
		}
		if s, ok := rawData["image_type"].(string); ok {
			result.BasicInfo.ImageType = s
		}
	}

	// content_safety
	if cs, ok := rawData["content_safety"].(map[string]interface{}); ok {
		var aiDecision bool
		var nsfwScore float64
		if v, ok := cs["is_nsfw"].(bool); ok {
			aiDecision = v
		}
		if v, ok := cs["nsfw_score"].(float64); ok {
			nsfwScore = v
			result.ContentSafety.NSFWScore = nsfwScore
		}
		result.ContentSafety.IsNSFW = applyNSFWThreshold(nsfwScore, aiDecision)
		if s, ok := cs["evaluation_result"].(string); ok {
			result.ContentSafety.EvaluationResult = s
		}
		if s, ok := cs["nsfw_reason"].(string); ok {
			result.ContentSafety.NSFWReason = s
		}
		if cat, ok := cs["categories"].(map[string]interface{}); ok {
			if v, ok := cat["nudity"].(float64); ok {
				result.ContentSafety.Categories.Nudity = v
			}
			if v, ok := cat["violence"].(float64); ok {
				result.ContentSafety.Categories.Violence = v
			}
			if v, ok := cat["hate_speech"].(float64); ok {
				result.ContentSafety.Categories.HateSpeech = v
			}
			if v, ok := cat["gambling"].(float64); ok {
				result.ContentSafety.Categories.Gambling = v
			}
			if v, ok := cat["alcohol_tobacco"].(float64); ok {
				result.ContentSafety.Categories.AlcoholTobacco = v
			}
		}
	} else if nsfw, ok := rawData["nsfw_detection"].(map[string]interface{}); ok {
		var aiDecision bool
		var nsfwScore float64
		if v, ok := nsfw["is_nsfw"].(bool); ok {
			aiDecision = v
		}
		if s, ok := nsfw["evaluation_result"].(string); ok {
			result.ContentSafety.EvaluationResult = s
		}
		if s, ok := nsfw["nsfw_reason"].(string); ok {
			result.ContentSafety.NSFWReason = s
		}
		if v, ok := nsfw["nudity_percentage"].(float64); ok {
			result.ContentSafety.Categories.Nudity = v / 100.0
		}
		if v, ok := nsfw["violence_percentage"].(float64); ok {
			result.ContentSafety.Categories.Violence = v / 100.0
		}
		if v, ok := nsfw["similarity_score"].(float64); ok {
			nsfwScore = v
			result.ContentSafety.NSFWScore = v
		}
		result.ContentSafety.IsNSFW = applyNSFWThreshold(nsfwScore, aiDecision)
	} else {
		result.ContentSafety.IsNSFW = false
		result.ContentSafety.NSFWScore = 0.0
		result.ContentSafety.EvaluationResult = "安全"
		result.ContentSafety.NSFWReason = ""
	}

	// 通用根字段提取（无论是否有 content_safety，都应解析这些字段）
	if s, ok := rawData["description"].(string); ok {
		result.Description = strings.TrimSpace(s)
	}
	if s, ok := rawData["search_content"].(string); ok {
		result.SearchContent = strings.TrimSpace(s)
	}
	if arr, ok := rawData["semantic_keywords"].([]interface{}); ok {
		kws := make([]string, 0, len(arr))
		for _, v := range arr {
			if sv, ok := v.(string); ok {
				sv = strings.TrimSpace(sv)
				if sv != "" {
					kws = append(kws, sv)
				}
			}
		}
		result.SemanticKeywords = kws
	}
	// is_recommended 兼容两种命名
	if v, ok := rawData["is_recommended"].(bool); ok {
		result.IsRecommended = v
	} else if v, ok := rawData["isRecommended"].(bool); ok {
		result.IsRecommended = v
	}
	if arr, ok := rawData["tags"].([]interface{}); ok {
		tags := make([]string, 0, len(arr))
		seen := make(map[string]struct{}, len(arr))
		for _, v := range arr {
			if sv, ok := v.(string); ok {
				sv = strings.TrimSpace(sv)
				if sv == "" {
					continue
				}
				key := strings.ToLower(sv)
				if _, dup := seen[key]; dup {
					continue
				}
				seen[key] = struct{}{}
				tags = append(tags, sv)
			}
		}
		result.Tags = tags
	}
	// visual_elements
	if ve, ok := rawData["visual_elements"].(map[string]interface{}); ok {
		if s, ok := ve["dominant_color"].(string); ok {
			result.VisualElements.DominantColor = s
		}
		if arr, ok := ve["color_palette"].([]interface{}); ok {
			colors := make([]string, 0, len(arr))
			for _, v := range arr {
				if sv, ok := v.(string); ok {
					sv = strings.TrimSpace(sv)
					if sv != "" {
						colors = append(colors, sv)
					}
				}
			}
			result.VisualElements.ColorPalette = colors
		}
		if s, ok := ve["composition"].(string); ok {
			result.VisualElements.Composition = s
		}
		if n, ok := ve["objects_count"].(float64); ok {
			result.VisualElements.ObjectsCount = int(n)
		}
	}

	return result, nil
}
