package vector

import (
	"context"
	"fmt"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// EmbeddingProvider 向量化提供者接口
type EmbeddingProvider interface {
	GenerateEmbedding(text string) ([]float32, error)
	BatchGenerateEmbeddings(texts []string) ([][]float32, error)
	GetDimension() int
	GetModel() string
}

// OpenAIEmbeddingClient OpenAI向量化客户端
type OpenAIEmbeddingClient struct {
	client    *openai.Client
	model     string
	dimension int
	timeout   time.Duration
}

// NewOpenAIEmbeddingClientWithConfig 使用指定配置创建OpenAI向量化客户端
func NewOpenAIEmbeddingClientWithConfig(apiKey, baseURL, model string, timeout int) *OpenAIEmbeddingClient {

	if apiKey == "" {
		logger.Warn("OpenAI API密钥未配置，向量化功能将不可用")
		return nil
	}

	if model == "" {
		model = string(openai.SmallEmbedding3) // 默认使用text-embedding-3-small
	}

	var dimension int
	switch model {
	case string(openai.LargeEmbedding3): // "text-embedding-3-large"
		dimension = 3072 // text-embedding-3-large 的维度
	case string(openai.SmallEmbedding3): // "text-embedding-3-small"
		dimension = 1536 // text-embedding-3-small 的维度
	default:
		logger.Warn("未识别的模型: %s，使用默认维度1536", model)
		dimension = 1536
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	if timeoutDuration == 0 {
		timeoutDuration = 30 * time.Second // 默认30秒超时
	}

	clientConfig := openai.DefaultConfig(apiKey)

	// 如果配置了自定义BaseURL，则使用代理地址（规范化URL，自动添加/v1）
	if baseURL != "" {
		clientConfig.BaseURL = utils.NormalizeOpenAIBaseURL(baseURL)
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIEmbeddingClient{
		client:    client,
		model:     model,
		dimension: dimension,
		timeout:   timeoutDuration,
	}
}

// GenerateEmbedding 生成单个文本的向量
func (c *OpenAIEmbeddingClient) GenerateEmbedding(text string) ([]float32, error) {
	if c == nil {
		return nil, fmt.Errorf("OpenAI客户端未初始化")
	}

	text = c.preprocessText(text)
	if text == "" {
		return nil, fmt.Errorf("文本内容为空")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(c.model),
	}

	resp, err := c.client.CreateEmbeddings(ctx, req)
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

// BatchGenerateEmbeddings 批量生成向量
func (c *OpenAIEmbeddingClient) BatchGenerateEmbeddings(texts []string) ([][]float32, error) {
	if c == nil {
		return nil, fmt.Errorf("OpenAI客户端未初始化")
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("文本列表为空")
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

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: processedTexts,
		Model: openai.EmbeddingModel(c.model),
	}

	resp, err := c.client.CreateEmbeddings(ctx, req)
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

func (c *OpenAIEmbeddingClient) GetDimension() int {
	if c == nil {
		return 0
	}
	return c.dimension
}

func (c *OpenAIEmbeddingClient) GetModel() string {
	if c == nil {
		return ""
	}
	return c.model
}

// preprocessText 预处理文本，优化向量化质量
func (c *OpenAIEmbeddingClient) preprocessText(text string) string {
	text = strings.TrimSpace(text)

	// 保留结构化信息，只替换多余的换行符
	// 保留句号后的换行符作为自然段落分隔
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

	// 动态限制文本长度（基于内容类型调整）
	maxLength := c.calculateOptimalLength(text)
	if len(text) > maxLength {
		text = c.smartTruncate(text, maxLength)
	}

	return text
}

// calculateOptimalLength 根据内容特点计算最佳长度
func (c *OpenAIEmbeddingClient) calculateOptimalLength(text string) int {
	baseLength := 6000

	// 如果包含结构化信息（如"主要内容："等标记），适当增加长度
	if strings.Contains(text, "主要内容：") || strings.Contains(text, "详细描述：") {
		baseLength = 8000
	}

	// 如果内容较短，不需要截断
	if len(text) < 1000 {
		return len(text) + 100
	}

	return baseLength
}

// smartTruncate 智能截断文本，保留完整语义
func (c *OpenAIEmbeddingClient) smartTruncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	truncated := text[:maxLength]
	if lastPeriod := strings.LastIndex(truncated, "。"); lastPeriod > maxLength-200 {
		return text[:lastPeriod+1]
	}

	// 尝试在句子结束符处截断
	sentenceEnders := []string{"！", "？", ".", "!", "?"}
	for _, ender := range sentenceEnders {
		if lastEnder := strings.LastIndex(truncated, ender); lastEnder > maxLength-200 {
			return text[:lastEnder+1]
		}
	}

	// 最后尝试在完整单词处截断
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLength-100 {
		return text[:lastSpace]
	}

	return truncated
}

// ValidateAPIKey 验证API密钥是否有效
func (c *OpenAIEmbeddingClient) ValidateAPIKey() error {
	if c == nil {
		return fmt.Errorf("客户端未初始化")
	}

	_, err := c.GenerateEmbedding("test")
	if err != nil {
		return fmt.Errorf("API密钥验证失败: %v", err)
	}

	return nil
}

// EstimateTokens 估算文本的token数量
func (c *OpenAIEmbeddingClient) EstimateTokens(text string) int {
	// 简单估算：英文约4字符/token，中文约1.5字符/token
	// 这是一个粗略估算，实际应该使用tiktoken库
	text = c.preprocessText(text)

	chineseChars := 0
	englishChars := 0

	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseChars++
		} else {
			englishChars++
		}
	}

	estimatedTokens := int(float64(chineseChars)/1.5 + float64(englishChars)/4.0)

	return estimatedTokens
}

// GetUsageCost 计算使用成本（美元）
func (c *OpenAIEmbeddingClient) GetUsageCost(tokenCount int) float64 {
	// text-embedding-3-small 价格: $0.02 / 1M tokens
	const pricePerMillion = 0.02
	return float64(tokenCount) * pricePerMillion / 1000000.0
}
