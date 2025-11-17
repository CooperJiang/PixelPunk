package ai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"pixelpunk/pkg/ai/prompts"
	"pixelpunk/pkg/logger"
	"strings"
	"time"
)

// OpenAIProvider OpenAI服务提供商实现
type OpenAIProvider struct {
	config *Config
	client *http.Client
}

func NewOpenAIProvider(config *Config) *OpenAIProvider {
	// 创建高性能HTTP客户端，支持高并发
	transport := &http.Transport{
		// 连接池优化 - 支持50个并发连接到OpenAI
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     90 * time.Second,

		// 禁用压缩以减少CPU开销（base64数据已压缩）
		DisableCompression:  false,
		DisableKeepAlives:   false,

		// 连接超时控制
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		// TLS优化
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},

		// 启用HTTP/2以支持多路复用
		ForceAttemptHTTP2: true,
	}

	client := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	return &OpenAIProvider{
		config: config,
		client: client,
	}
}

// AnalyzeFile 分析文件 - 支持URL和base64数据
func (p *OpenAIProvider) AnalyzeFile(ctx context.Context, req *FileAnalysisRequest) (*AIResponse, error) {
	var imageContent map[string]interface{}

	if req.ImageData != "" {
		imageContent = map[string]interface{}{
			"url": fmt.Sprintf("data:image/%s;base64,%s", req.Format, req.ImageData),
		}
	} else if req.ImageURL != "" {
		imageContent = map[string]interface{}{
			"url": req.ImageURL,
		}
	} else {
		return &AIResponse{
			Success: false,
			ErrMsg:  "缺少文件数据或URL",
		}, nil
	}

	// 用户文本：为空时使用默认分析提示词
	userText := strings.TrimSpace(req.Prompt)
	if userText == "" {
		userText = prompts.GetFileAnalysisPrompt()
	}

	requestMap := map[string]interface{}{
		"model":       p.config.Model,
		"max_tokens":  p.config.MaxTokens,
		"temperature": p.config.Temperature,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": getImageAnalysisSystemPrompt(),
			},
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": userText,
					},
					{
						"type":      "image_url",
						"image_url": imageContent,
					},
				},
			},
		},
	}

	return p.sendRequest(ctx, requestMap, req.ImageURL)
}

// CategorizeFile 文件分类
func (p *OpenAIProvider) CategorizeFile(ctx context.Context, req *FileCategorizationRequest) (*FileCategorizationResponse, error) {
	if len(req.Categories) == 0 {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "没有可用的分类选项",
		}, nil
	}

	var imageContent map[string]interface{}

	if req.ImageData != "" {
		imageContent = map[string]interface{}{
			"url": fmt.Sprintf("data:image/%s;base64,%s", req.Format, req.ImageData),
		}
	} else if req.ImageURL != "" {
		imageContent = map[string]interface{}{
			"url": req.ImageURL,
		}
	} else {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "缺少文件数据或URL",
		}, nil
	}

	promptCategories := make([]prompts.CategoryInfo, len(req.Categories))
	for i, cat := range req.Categories {
		promptCategories[i] = prompts.CategoryInfo{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Source:      cat.Source,
		}
	}

	prompt := prompts.GetFileCategorizationPrompt(promptCategories)
	systemPrompt := prompts.GetFileCategorizationSystemPrompt()

	requestMap := map[string]interface{}{
		"model":       p.config.Model,
		"max_tokens":  p.config.MaxTokens,
		"temperature": 0.1, // 分类任务使用较低的temperature确保一致性
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type":      "image_url",
						"image_url": imageContent,
					},
				},
			},
		},
	}

	response, err := p.sendRequest(ctx, requestMap, req.ImageURL)
	if err != nil {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI分类请求失败: %v", err),
		}, err
	}

	if !response.Success {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  response.ErrMsg,
		}, nil
	}

	return p.parseCategorizationResponse(response.Data, response.Usage)
}

// parseCategorizationResponse 解析分类响应
func (p *OpenAIProvider) parseCategorizationResponse(content string, usage *TokenUsage) (*FileCategorizationResponse, error) {
	var result struct {
		Success             bool   `json:"success"`
		CategoryID          uint   `json:"category_id"`
		CategoryName        string `json:"category_name"`
		CategoryDescription string `json:"category_description"`
	}

	cleanContent := CleanJSON(ExtractJSONFromText(content))

	err := json.Unmarshal([]byte(cleanContent), &result)
	if err != nil {
		logger.Error("解析AI分类响应失败: %v, content: %s", err, cleanContent)
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("解析AI分类结果失败: %v", err),
		}, nil
	}

	if !result.Success {
		logger.Warn("AI分类失败: Success=%t, 原始内容: %s", result.Success, cleanContent)
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "AI未能成功识别文件分类",
		}, nil
	}

	// CategoryID=0 表示AI建议创建新分类，这是正常情况

	return &FileCategorizationResponse{
		Success:             true,
		CategoryID:          result.CategoryID,
		CategoryName:        result.CategoryName,
		CategoryDescription: result.CategoryDescription,
		Usage:               usage,
	}, nil
}

// GenerateEmbedding 生成文本向量
func (p *OpenAIProvider) GenerateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	if strings.TrimSpace(req.Text) == "" {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "文本内容不能为空",
		}, nil
	}

	// 使用指定模型或默认的文本向量化模型
	model := req.Model
	if model == "" {
		model = "text-embedding-3-small" // 默认向量化模型
	}

	requestMap := map[string]interface{}{
		"model": model,
		"input": req.Text,
	}

	return p.sendEmbeddingRequest(ctx, requestMap, model)
}

// TestConnection 测试连接
func (p *OpenAIProvider) TestConnection(ctx context.Context) (*TestResult, error) {
	startTime := time.Now()

	requestMap := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "请回复：连接正常",
			},
		},
		"max_tokens":  10,
		"temperature": 0,
	}

	response, err := p.sendRequest(ctx, requestMap, "")
	responseTime := time.Since(startTime)

	result := &TestResult{
		ResponseTime: responseTime,
	}

	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("连接测试失败: %v", err)
		return result, nil
	}

	if response.Success {
		result.Success = true
		result.Message = "AI配置测试成功"
		result.Details = map[string]interface{}{
			"model":         p.config.Model,
			"api_endpoint":  p.getAPIURL(),
			"test_response": response.Data,
			"status":        "connected",
		}

		if response.Usage != nil {
			result.Details["tokens_used"] = response.Usage.TotalTokens
		}
	} else {
		result.Success = false
		result.Message = response.ErrMsg
	}

	return result, nil
}

func (p *OpenAIProvider) GetProviderInfo() *ProviderInfo {
	return &ProviderInfo{
		Name:        "openai",
		DisplayName: "OpenAI",
		Models:      []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-4-vision-preview"},
		Features:    []string{"text_generation", "image_analysis", "vision"},
	}
}

// sendRequest 发送HTTP请求到OpenAI
func (p *OpenAIProvider) sendRequest(ctx context.Context, requestMap map[string]interface{}, imageURL string) (*AIResponse, error) {
	requestJSON, err := json.Marshal(requestMap)
	if err != nil {
		logger.Error("序列化OpenAI请求失败: %v", err)
		return &AIResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("序列化请求失败: %v", err),
		}, nil
	}

	apiURL := p.getAPIURL()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		logger.Error("创建HTTP请求失败: %v", err)
		return &AIResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("创建HTTP请求失败: %v", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	// 记录HTTP调用开始时间
	httpStart := time.Now()
	resp, err := p.client.Do(req)
	httpDuration := time.Since(httpStart).Milliseconds()
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") ||
			strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "Client.Timeout") {
			logger.Error("请求OpenAI超时: %v", err)
			return &AIResponse{
				Success: false,
				ErrMsg:  fmt.Sprintf("请求超时，请稍后重试: %v", err),
			}, nil
		}

		logger.Error("发送请求到OpenAI失败: %v", err)
		return &AIResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("发送请求失败: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取OpenAI响应失败: %v", err)
		return &AIResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("读取响应失败: %v", err),
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return p.handleErrorResponseWithDuration(resp.StatusCode, respBody, apiURL, httpDuration)
	}

	return p.parseSuccessResponse(respBody, imageURL, httpDuration)
}

// handleErrorResponseWithDuration 处理错误响应（带HTTP耗时）
func (p *OpenAIProvider) handleErrorResponseWithDuration(statusCode int, respBody []byte, apiURL string, httpDuration int64) (*AIResponse, error) {
	resp, err := p.handleErrorResponse(statusCode, respBody, apiURL)
	if resp != nil {
		resp.HttpDuration = httpDuration
	}
	return resp, err
}

// handleErrorResponse 处理错误响应
func (p *OpenAIProvider) handleErrorResponse(statusCode int, respBody []byte, apiURL string) (*AIResponse, error) {
	var errorResp struct {
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	errorMsg := fmt.Sprintf("API调用失败, 状态码: %d", statusCode)

	if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error.Message != "" {
		errorMsg = p.sanitizeErrorMessage(errorResp.Error.Message)
	}

	// 针对常见错误码提供更友好的提示
	var friendlyMessage string
	switch statusCode {
	case 401:
		friendlyMessage = "API密钥无效或未设置"
	case 403:
		friendlyMessage = "API密钥没有权限访问此模型"
	case 404:
		friendlyMessage = fmt.Sprintf("API端点不存在，请检查代理地址: %s", apiURL)
	case 429:
		friendlyMessage = "API请求频率超限，请稍后重试"
	case 500:
		friendlyMessage = "OpenAI服务器内部错误"
	case 502:
		friendlyMessage = "网关错误，请检查代理服务器"
	case 503:
		friendlyMessage = "服务暂时不可用"
	default:
		friendlyMessage = errorMsg
	}

	logger.Error("OpenAI API错误: %s", friendlyMessage)
	return &AIResponse{
		Success: false,
		ErrMsg:  friendlyMessage,
	}, nil
}

// parseSuccessResponse 解析成功响应
func (p *OpenAIProvider) parseSuccessResponse(respBody []byte, imageURL string, httpDuration int64) (*AIResponse, error) {
	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage *TokenUsage `json:"usage"`
	}

	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		logger.Error("解析OpenAI响应失败: %v", err)
		return &AIResponse{
			Success:      false,
			ErrMsg:       fmt.Sprintf("解析响应失败: %v", err),
			HttpDuration: httpDuration,
		}, nil
	}

	if len(openaiResp.Choices) == 0 {
		return &AIResponse{
			Success:      false,
			ErrMsg:       "OpenAI返回空响应",
			HttpDuration: httpDuration,
		}, nil
	}

	content := strings.TrimSpace(openaiResp.Choices[0].Message.Content)
	if content == "" {
		return &AIResponse{
			Success:      false,
			ErrMsg:       "OpenAI返回空内容",
			HttpDuration: httpDuration,
		}, nil
	}

	return &AIResponse{
		Success:      true,
		Data:         content,
		ImageURL:     imageURL,
		Usage:        openaiResp.Usage,
		HttpDuration: httpDuration,
	}, nil
}

// sendEmbeddingRequest 发送向量化请求到OpenAI
func (p *OpenAIProvider) sendEmbeddingRequest(ctx context.Context, requestMap map[string]interface{}, model string) (*EmbeddingResponse, error) {
	requestJSON, err := json.Marshal(requestMap)
	if err != nil {
		logger.Error("序列化OpenAI向量化请求失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("序列化请求失败: %v", err),
		}, nil
	}

	apiURL := p.getEmbeddingAPIURL()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		logger.Error("创建HTTP请求失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("创建HTTP请求失败: %v", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") ||
			strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "Client.Timeout") {
			logger.Error("请求OpenAI向量化超时: %v", err)
			return &EmbeddingResponse{
				Success: false,
				ErrMsg:  fmt.Sprintf("请求超时，请稍后重试: %v", err),
			}, nil
		}

		logger.Error("发送向量化请求到OpenAI失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("发送请求失败: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取OpenAI向量化响应失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("读取响应失败: %v", err),
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return p.handleEmbeddingErrorResponse(resp.StatusCode, respBody, apiURL)
	}

	return p.parseEmbeddingResponse(respBody, model)
}

// getAPIURL 获取聊天API URL
func (p *OpenAIProvider) getAPIURL() string {
	base := strings.TrimRight(strings.TrimSpace(p.config.BaseURL), "/")
	// If base already ends with /v1, don't append it twice
	if strings.HasSuffix(base, "/v1") {
		return base + "/chat/completions"
	}
	return base + "/v1/chat/completions"
}

// getEmbeddingAPIURL 获取向量化API URL
func (p *OpenAIProvider) getEmbeddingAPIURL() string {
	base := strings.TrimRight(strings.TrimSpace(p.config.BaseURL), "/")
	if strings.HasSuffix(base, "/v1") {
		return base + "/embeddings"
	}
	return base + "/v1/embeddings"
}

// handleEmbeddingErrorResponse 处理向量化错误响应
func (p *OpenAIProvider) handleEmbeddingErrorResponse(statusCode int, respBody []byte, apiURL string) (*EmbeddingResponse, error) {
	var errorResp struct {
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	errorMsg := fmt.Sprintf("向量化API调用失败, 状态码: %d", statusCode)

	if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error.Message != "" {
		errorMsg = p.sanitizeErrorMessage(errorResp.Error.Message)
	}

	// 针对常见错误码提供更友好的提示
	var friendlyMessage string
	switch statusCode {
	case 401:
		friendlyMessage = "API密钥无效或未设置"
	case 403:
		friendlyMessage = "API密钥没有权限访问向量化模型"
	case 404:
		friendlyMessage = fmt.Sprintf("向量化API端点不存在，请检查代理地址: %s", apiURL)
	case 429:
		friendlyMessage = "API请求频率超限，请稍后重试"
	case 500:
		friendlyMessage = "OpenAI向量化服务器内部错误"
	case 502:
		friendlyMessage = "网关错误，请检查代理服务器"
	case 503:
		friendlyMessage = "向量化服务暂时不可用"
	default:
		friendlyMessage = errorMsg
	}

	logger.Error("OpenAI向量化API错误: %s", friendlyMessage)
	return &EmbeddingResponse{
		Success: false,
		ErrMsg:  friendlyMessage,
	}, nil
}

// parseEmbeddingResponse 解析向量化成功响应
func (p *OpenAIProvider) parseEmbeddingResponse(respBody []byte, model string) (*EmbeddingResponse, error) {
	var openaiResp struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
			Object    string    `json:"object"`
		} `json:"data"`
		Model string      `json:"model"`
		Usage *TokenUsage `json:"usage"`
	}

	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		logger.Error("解析OpenAI向量化响应失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("解析响应失败: %v", err),
		}, nil
	}

	if len(openaiResp.Data) == 0 {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "OpenAI返回空向量响应",
		}, nil
	}

	// 转换float64数组为float32数组
	embedding64 := openaiResp.Data[0].Embedding
	embedding32 := make([]float32, len(embedding64))
	for i, v := range embedding64 {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			logger.Warn("向量中包含无效值，位置: %d, 值: %f", i, v)
			embedding32[i] = 0.0
		} else {
			embedding32[i] = float32(v)
		}
	}

	dimension := len(embedding32)
	if dimension == 0 {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "OpenAI返回空向量",
		}, nil
	}

	return &EmbeddingResponse{
		Success:   true,
		Embedding: embedding32,
		Model:     model,
		Dimension: dimension,
		Usage:     openaiResp.Usage,
	}, nil
}

// sanitizeErrorMessage 清理错误消息，移除敏感信息
func (p *OpenAIProvider) sanitizeErrorMessage(message string) string {
	message = strings.ReplaceAll(message, p.config.APIKey, "[API_KEY]")
	return message
}

// getImageAnalysisSystemPrompt 获取文件分析系统提示词
func getImageAnalysisSystemPrompt() string {
	// 使用 prompts 包中的完整系统提示词，包含详细的NSFW检测规则
	return prompts.GetFileAnalysisSystemPrompt()
}

// fixJSONTrailingCommas 已迁移至 jsonutil.go

// TagFile 文件标注（支持标签列表）
func (p *OpenAIProvider) TagFile(ctx context.Context, req *FileTaggingRequest) (*FileAnalysisResponse, error) {
	// 构建专门用于标签生成的提示词（改为JSON返回，仍兼容老格式）
	var ptags []prompts.TagInfo
	if len(req.AvailableTags) > 0 {
		ptags = make([]prompts.TagInfo, 0, len(req.AvailableTags))
		for _, t := range req.AvailableTags {
			ptags = append(ptags, prompts.TagInfo{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
				Source:      t.Source,
				UsageCount:  t.UsageCount,
			})
		}
	}
	promptText := prompts.GetFileTaggingPrompt(ptags, true)

	var imageContent map[string]interface{}

	if req.ImageData != "" {
		imageContent = map[string]interface{}{
			"url": fmt.Sprintf("data:image/%s;base64,%s", req.Format, req.ImageData),
		}
	} else if req.ImageURL != "" {
		imageContent = map[string]interface{}{
			"url": req.ImageURL,
		}
	} else {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  "必须提供文件URL或base64数据",
		}, nil
	}

	requestMap := map[string]interface{}{
		"model":       p.config.Model,
		"max_tokens":  p.config.MaxTokens,
		"temperature": p.config.Temperature,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": promptText,
					},
					{
						"type":      "image_url",
						"image_url": imageContent,
					},
				},
			},
		},
	}

	response, err := p.sendRequest(ctx, requestMap, "")
	if err != nil {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI请求失败: %v", err),
		}, err
	}

	if !response.Success {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  response.ErrMsg,
		}, nil
	}

	// 优先解析JSON；失败则回退到逗号分割
	raw := strings.TrimSpace(response.Data)
	cleaned := CleanJSON(ExtractJSONFromText(raw))
	var jsonResult struct {
		Tags        []string `json:"tags"`
		Description string   `json:"description"`
	}
	var tags []string
	desc := ""
	if cleaned != "" && json.Unmarshal([]byte(cleaned), &jsonResult) == nil && len(jsonResult.Tags) > 0 {
		// 去重/清理空白
		uniq := make(map[string]struct{}, len(jsonResult.Tags))
		for _, t := range jsonResult.Tags {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			if _, ok := uniq[strings.ToLower(t)]; ok {
				continue
			}
			uniq[strings.ToLower(t)] = struct{}{}
			tags = append(tags, t)
		}
		desc = strings.TrimSpace(jsonResult.Description)
	} else {
		// 兼容旧格式：逗号分割，并做去重与长度限制
		content := strings.ReplaceAll(raw, "，", ",")
		uniq := make(map[string]struct{})
		for _, t := range strings.Split(content, ",") {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			key := strings.ToLower(t)
			if _, ok := uniq[key]; ok {
				continue
			}
			uniq[key] = struct{}{}
			tags = append(tags, t)
		}
		desc = strings.TrimSpace(raw)
	}

	// 数量严格控制在最多7个（提示词已约束，这里做兜底）
	if len(tags) > 7 {
		tags = tags[:7]
	}

	return &FileAnalysisResponse{
		Success:     true,
		Tags:        tags,
		Description: desc,
		Usage:       response.Usage,
	}, nil
}
