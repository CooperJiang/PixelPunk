package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// DynamicOpenAIClient 动态OpenAI客户端（无需初始化，每次调用时读取最新配置）
type DynamicOpenAIClient struct{}

func NewDynamicOpenAIClient() *DynamicOpenAIClient {
	return &DynamicOpenAIClient{}
}

// getConfigFromDB 从数据库动态读取配置（绕过缓存，直接查数据库）
func (c *DynamicOpenAIClient) getConfigFromDB() (apiKey, baseURL, model string, timeout time.Duration, dimension int, err error) {
	// 直接从数据库读取配置（不使用缓存）
	// 原因：缓存可能在启动时保存了空配置，导致后续配置更新无法生效
	db := database.GetDB()
	if db == nil {
		return "", "", "", 0, 0, fmt.Errorf("数据库连接不可用")
	}

	type SettingRow struct {
		Key   string
		Value string
		Type  string
	}

	var settings []SettingRow
	if err := db.Table("setting").
		Where("`group` = ?", "vector").
		Where("`key` IN (?)", []string{"vector_api_key", "vector_base_url", "vector_model", "vector_timeout"}).
		Select("`key`, `value`, `type`").
		Find(&settings).Error; err != nil {
		return "", "", "", 0, 0, fmt.Errorf("查询向量配置失败: %v", err)
	}

	apiKey = ""
	baseURL = "https://api.openai.com/v1"
	model = "text-embedding-3-small"
	timeout = 30 * time.Second
	dimension = 1536

	for _, s := range settings {
		switch s.Key {
		case "vector_api_key":
			// 字符串需要JSON解析（数据库存储的是JSON格式）
			var key string
			if err := json.Unmarshal([]byte(s.Value), &key); err == nil {
				apiKey = key
			}
		case "vector_base_url":
			// 字符串需要JSON解析
			var url string
			if err := json.Unmarshal([]byte(s.Value), &url); err == nil && url != "" {
				// 标准化URL：自动添加/v1，去除多余斜杠等
				baseURL = utils.NormalizeOpenAIBaseURL(url)
			}
		case "vector_model":
			// 字符串需要JSON解析
			var modelName string
			if err := json.Unmarshal([]byte(s.Value), &modelName); err == nil && modelName != "" {
				model = modelName
			}
		case "vector_timeout":
			var timeoutSeconds float64
			if err := json.Unmarshal([]byte(s.Value), &timeoutSeconds); err == nil && timeoutSeconds > 0 {
				timeout = time.Duration(timeoutSeconds) * time.Second
			}
		}
	}

	switch model {
	case "text-embedding-3-large":
		dimension = 3072
	case "text-embedding-3-small":
		dimension = 1536
	default:
		dimension = 1536
	}

	if apiKey == "" {
		return "", "", "", 0, 0, fmt.Errorf("OpenAI API密钥未配置")
	}

	return apiKey, baseURL, model, timeout, dimension, nil
}

// GenerateEmbedding 生成单个文本的向量（动态读取配置）
func (c *DynamicOpenAIClient) GenerateEmbedding(text string) ([]float32, error) {
	apiKey, baseURL, model, timeout, _, err := c.getConfigFromDB()
	if err != nil {
		return nil, fmt.Errorf("读取向量配置失败: %v", err)
	}

	text = c.preprocessText(text)
	if text == "" {
		return nil, fmt.Errorf("文本内容为空")
	}

	clientConfig := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		clientConfig.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(clientConfig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(model),
	}

	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		logger.Error("OpenAI API调用失败: %v", err)
		return nil, fmt.Errorf("OpenAI向量化失败: %v", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("OpenAI返回空向量数据")
	}

	embedding := resp.Data[0].Embedding

	vector := make([]float32, len(embedding))
	for i, v := range embedding {
		vector[i] = float32(v)
	}

	return vector, nil
}

// BatchGenerateEmbeddings 批量生成向量（动态读取配置）
func (c *DynamicOpenAIClient) BatchGenerateEmbeddings(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("文本列表为空")
	}

	apiKey, baseURL, model, timeout, _, err := c.getConfigFromDB()
	if err != nil {
		return nil, fmt.Errorf("读取向量配置失败: %v", err)
	}

	processedTexts := make([]string, 0, len(texts))
	for _, text := range texts {
		processed := c.preprocessText(text)
		if processed != "" {
			processedTexts = append(processedTexts, processed)
		}
	}

	if len(processedTexts) == 0 {
		return nil, fmt.Errorf("没有有效的文本内容")
	}

	clientConfig := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		clientConfig.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(clientConfig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: processedTexts,
		Model: openai.EmbeddingModel(model),
	}

	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		logger.Error("OpenAI批量向量化请求失败: %v", err)
		return nil, fmt.Errorf("OpenAI批量向量化失败: %v", err)
	}

	if len(resp.Data) != len(processedTexts) {
		return nil, fmt.Errorf("返回向量数量不匹配，期望: %d, 实际: %d", len(processedTexts), len(resp.Data))
	}

	vectors := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		vector := make([]float32, len(data.Embedding))
		for j, v := range data.Embedding {
			vector[j] = float32(v)
		}
		vectors[i] = vector
	}

	return vectors, nil
}

func (c *DynamicOpenAIClient) GetDimension() int {
	_, _, _, _, dimension, err := c.getConfigFromDB()
	if err != nil {
		logger.Warn("获取向量维度失败，使用默认值1536: %v", err)
		return 1536
	}
	return dimension
}

func (c *DynamicOpenAIClient) GetModel() string {
	_, _, model, _, _, err := c.getConfigFromDB()
	if err != nil {
		logger.Warn("获取向量模型失败，使用默认值: %v", err)
		return "text-embedding-3-small"
	}
	return model
}

// preprocessText 预处理文本，优化向量化质量
func (c *DynamicOpenAIClient) preprocessText(text string) string {
	text = strings.TrimSpace(text)

	// 保留结构化信息，只替换多余的换行符
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 将多个连续换行符合并为双换行符（保留段落结构）
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	// 将单个换行符替换为空格（保持句子连贯性）
	lines := strings.Split(text, "\n\n")
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(strings.TrimSpace(line), "\n", " ")
	}
	text = strings.Join(lines, "\n\n")

	text = strings.ReplaceAll(text, "\t", " ")

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// 优化中文标点符号后的空格处理
	text = strings.ReplaceAll(text, "。 ", "。")
	text = strings.ReplaceAll(text, "， ", "，")
	text = strings.ReplaceAll(text, "！ ", "！")
	text = strings.ReplaceAll(text, "？ ", "？")

	maxLength := c.calculateOptimalLength(text)
	if len(text) > maxLength {
		text = c.smartTruncate(text, maxLength)
	}

	return text
}

// calculateOptimalLength 根据内容特点计算最佳长度
func (c *DynamicOpenAIClient) calculateOptimalLength(text string) int {
	baseLength := 6000

	if strings.Contains(text, "主要内容：") || strings.Contains(text, "详细描述：") {
		baseLength = 8000
	}

	if len(text) < 1000 {
		return len(text) + 100
	}

	return baseLength
}

// smartTruncate 智能截断文本，保留完整语义
func (c *DynamicOpenAIClient) smartTruncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	truncated := text[:maxLength]
	if lastPeriod := strings.LastIndex(truncated, "。"); lastPeriod > maxLength-200 {
		return text[:lastPeriod+1]
	}

	sentenceEnders := []string{"！", "？", ".", "!", "?"}
	for _, ender := range sentenceEnders {
		if lastEnder := strings.LastIndex(truncated, ender); lastEnder > maxLength-200 {
			return text[:lastEnder+1]
		}
	}

	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLength-100 {
		return text[:lastSpace]
	}

	return truncated
}
