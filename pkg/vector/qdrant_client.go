package vector

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
)

// QdrantClient 直接连接Qdrant的简单客户端
type QdrantClient struct {
	baseURL    string
	httpClient *http.Client
	collection string
}

// QdrantPoint Qdrant点结构
type QdrantPoint struct {
	Id      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

// QdrantSearchRequest 搜索请求
type QdrantSearchRequest struct {
	Vector         []float32              `json:"vector"`
	Filter         map[string]interface{} `json:"filter,omitempty"`
	Limit          int                    `json:"limit"`
	WithPayload    bool                   `json:"with_payload"`
	WithVector     bool                   `json:"with_vector"`
	ScoreThreshold *float32               `json:"score_threshold,omitempty"`
}

// QdrantSearchResponse 搜索响应
type QdrantSearchResponse struct {
	Result []QdrantSearchResult `json:"result"`
	Status string               `json:"status"`
	Time   float64              `json:"time"`
}

// QdrantSearchResult 搜索结果
type QdrantSearchResult struct {
	Id      string                 `json:"id"`
	Version int                    `json:"version"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
	Vector  []float32              `json:"vector,omitempty"`
}

func NewQdrantClient(qdrantURL string, timeout int) *QdrantClient {
	return &QdrantClient{
		baseURL:    qdrantURL,
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
		collection: "file_vectors",
	}
}

// generateQdrantID 基于文件ID生成确定性UUID
func (q *QdrantClient) generateQdrantID(fileID string) string {
	// 使用MD5哈希生成确定性的UUID
	hash := md5.Sum([]byte(fileID))
	hashHex := hex.EncodeToString(hash[:])

	// 将32位十六进制字符串格式化为UUID格式
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		hashHex[0:8],
		hashHex[8:12],
		hashHex[12:16],
		hashHex[16:20],
		hashHex[20:32])

	return uuidStr
}

// InitCollection 初始化向量集合
func (q *QdrantClient) InitCollection() error {
	resp, err := q.httpClient.Get(fmt.Sprintf("%s/collections/%s", q.baseURL, q.collection))
	if err == nil && resp != nil && resp.StatusCode == 200 {
		resp.Body.Close()
		return nil
	}
	if resp != nil {
		resp.Body.Close()
	}

	createReq := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536, // text-embedding-3-small 向量维度
			"distance": "Cosine",
		},
	}

	reqBody, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("序列化创建集合请求失败: %w", err)
	}

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/collections/%s", q.baseURL, q.collection),
		bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建PUT请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = q.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("创建集合请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建集合失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ProcessFile 处理文件向量化
func (q *QdrantClient) ProcessFile(fileID, description string, userID uint) error {
	// 这里需要调用向量化服务生成向量
	return fmt.Errorf("向量化功能需要集成OpenAI客户端")
}

// StoreVector 存储向量
func (q *QdrantClient) StoreVector(fileID string, vector []float32, description string, model string) error {
	// 生成基于 fileID 的确定性 UUID
	qdrantID := q.generateQdrantID(fileID)

	// 可选：查询用户ID以写入payload，便于按用户过滤
	var userID uint
	if db := database.GetDB(); db != nil {
		var file models.File
		if err := db.Where("id = ?", fileID).First(&file).Error; err == nil {
			userID = file.UserID
		}
	}

	point := QdrantPoint{
		Id:     qdrantID,
		Vector: vector,
		Payload: map[string]interface{}{
			"file_id":     fileID, // 保存原始文件ID
			"description": description,
			"model":       model,
			"user_id":     userID,
		},
	}

	reqBody := map[string]interface{}{
		"points": []QdrantPoint{point},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化存储请求失败: %w", err)
	}

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/collections/%s/points", q.baseURL, q.collection),
		bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("创建PUT请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("存储向量请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("存储向量失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchVectors 搜索相似向量
func (q *QdrantClient) SearchVectors(queryVector []float32, limit int, userID uint, threshold float32) ([]VectorSearchResult, error) {
	searchReq := QdrantSearchRequest{
		Vector:         queryVector,
		Limit:          limit,
		WithPayload:    true,
		WithVector:     false,
		ScoreThreshold: &threshold,
	}

	if userID > 0 {
		searchReq.Filter = map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"key":   "user_id",
					"match": map[string]interface{}{"value": userID},
				},
			},
		}
	}

	reqData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("序列化搜索请求失败: %w", err)
	}

	resp, err := q.httpClient.Post(
		fmt.Sprintf("%s/collections/%s/points/search", q.baseURL, q.collection),
		"application/json",
		bytes.NewBuffer(reqData),
	)
	if err != nil {
		return nil, fmt.Errorf("搜索请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("搜索失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var searchResp QdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("解析搜索响应失败: %w", err)
	}

	results := make([]VectorSearchResult, 0, len(searchResp.Result))
	for _, item := range searchResp.Result {
		result := VectorSearchResult{
			Score:      item.Score,
			Similarity: item.Score, // 相似度就是Qdrant返回的Score
		}

		// 从 payload 中提取原始文件ID
		if fileID, ok := item.Payload["file_id"].(string); ok {
			result.FileID = fileID
		} else {
			// 如果payload中没有file_id，跳过这个结果
			logger.Warn("搜索结果中缺少file_id，跳过: %s", item.Id)
			continue
		}

		// 从payload中提取描述信息
		if desc, ok := item.Payload["description"].(string); ok {
			result.Description = desc
		}

		results = append(results, result)
	}

	return results, nil
}

// HealthCheck 健康检查
func (q *QdrantClient) HealthCheck() error {
	resp, err := q.httpClient.Get(q.baseURL + "/")
	if err != nil {
		return fmt.Errorf("qdrant健康检查失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("qdrant服务不健康，状态码: %d", resp.StatusCode)
	}

	return nil
}

// 当前未实现，返回0（该功能暂不需要）
func (q *QdrantClient) GetUserVectorCount(userID uint) (int64, error) {
	return 0, nil
}

// GetAllFileIDs 从 Qdrant 遍历获取 file_id（通过 payload.file_id），用于对账/清理孤儿
// limit: 最多返回的数量；<=0 表示不限制（但为安全起见这里按批次滚动，最多返回 100k）
func (q *QdrantClient) GetAllFileIDs(limit int) ([]string, error) {
	type scrollReq struct {
		WithPayload bool                    `json:"with_payload"`
		WithVector  bool                    `json:"with_vector"`
		Limit       int                     `json:"limit"`
		Offset      *map[string]interface{} `json:"offset,omitempty"`
	}
	type point struct {
		ID      interface{}            `json:"id"`
		Payload map[string]interface{} `json:"payload"`
	}
	type scrollResp struct {
		Result struct {
			Points []point                 `json:"points"`
			Offset *map[string]interface{} `json:"offset"`
		} `json:"result"`
		Status string  `json:"status"`
		Time   float64 `json:"time"`
	}

	// 滚动读取，每批 1000
	batch := 1000
	if limit > 0 && limit < batch {
		batch = limit
	}
	maxTotal := 100000
	if limit > 0 && limit < maxTotal {
		maxTotal = limit
	}

	var all []string
	var offset *map[string]interface{}
	for {
		reqBody := scrollReq{WithPayload: true, WithVector: false, Limit: batch, Offset: offset}
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("序列化scroll请求失败: %w", err)
		}

		url := fmt.Sprintf("%s/collections/%s/points/scroll", q.baseURL, q.collection)
		resp, err := q.httpClient.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("scroll 请求失败: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			b, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("scroll 失败，状态码: %d, 响应: %s", resp.StatusCode, string(b))
		}

		var sr scrollResp
		if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
			return nil, fmt.Errorf("解析scroll响应失败: %w", err)
		}
		if sr.Status != "ok" {
			return nil, fmt.Errorf("scroll 响应状态异常: %s", sr.Status)
		}

		for _, p := range sr.Result.Points {
			if p.Payload != nil {
				if v, ok := p.Payload["file_id"].(string); ok && v != "" {
					all = append(all, v)
				}
			}
		}

		if len(all) >= maxTotal {
			break
		}
		if sr.Result.Offset == nil {
			break
		}
		offset = sr.Result.Offset
	}

	// 截断到 limit
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func (q *QdrantClient) GetStorageStats() (*VectorStorageStats, error) {
	// 优先使用 /points/count 以提升跨版本兼容性
	type countResp struct {
		Result struct {
			Count int64 `json:"count"`
		} `json:"result"`
		Status string `json:"status"`
	}

	countURL := fmt.Sprintf("%s/collections/%s/points/count", q.baseURL, q.collection)
	reqBody := bytes.NewBuffer([]byte(`{"exact": true}`))
	resp, err := q.httpClient.Post(countURL, "application/json", reqBody)
	if err == nil && resp != nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var cr countResp
		if err := json.NewDecoder(resp.Body).Decode(&cr); err == nil && cr.Status == "ok" {
			return &VectorStorageStats{
				TotalVectors:   cr.Result.Count,
				CompletedCount: cr.Result.Count,
				PendingCount:   0,
				FailedCount:    0,
				AverageNorm:    0,
				LastUpdateTime: time.Now(),
				StorageSize:    0,
			}, nil
		}
	}
	if resp != nil {
		resp.Body.Close()
	}

	// 回退到 /collections/<name>
	resp2, err2 := q.httpClient.Get(fmt.Sprintf("%s/collections/%s", q.baseURL, q.collection))
	if err2 != nil {
		return nil, fmt.Errorf("获取集合信息失败: %w", err2)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		body, _ := io.ReadAll(resp2.Body)
		return nil, fmt.Errorf("获取集合信息失败，状态码: %d, 响应: %s", resp2.StatusCode, string(body))
	}

	var collectionResp struct {
		Result struct {
			Status              string `json:"status"`
			IndexedVectorsCount int64  `json:"indexed_vectors_count"`
			PointsCount         int64  `json:"points_count"`
		} `json:"result"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&collectionResp); err != nil {
		return nil, fmt.Errorf("解析集合信息失败: %w", err)
	}
	if collectionResp.Status != "ok" {
		return nil, fmt.Errorf("集合信息响应状态异常: %s", collectionResp.Status)
	}
	return &VectorStorageStats{
		TotalVectors:   collectionResp.Result.PointsCount,
		CompletedCount: collectionResp.Result.IndexedVectorsCount,
		PendingCount:   0,
		FailedCount:    0,
		AverageNorm:    0,
		LastUpdateTime: time.Now(),
		StorageSize:    0,
	}, nil
}

// VectorExists 检查向量是否存在
func (q *QdrantClient) VectorExists(fileID string) (bool, error) {
	// 生成 Qdrant 点ID
	pointID := q.generateQdrantID(fileID)

	// 构造查询向量点的URL
	url := fmt.Sprintf("%s/collections/%s/points/%s", q.baseURL, q.collection, pointID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("qdrant 响应错误 %d: %s", resp.StatusCode, string(body))
	}
}

// BatchProcessFiles 批量处理文件
func (q *QdrantClient) BatchProcessFiles(items []VectorItem) error {
	for _, item := range items {
		// Since VectorItem doesn't have UserID, we'll use 0 as default or handle it differently
		if err := q.ProcessFile(item.FileID, item.Description, 0); err != nil {
			return fmt.Errorf("批量处理失败，文件ID: %s, 错误: %w", item.FileID, err)
		}
	}
	return nil
}

// BatchStoreVectors 批量存储向量
func (q *QdrantClient) BatchStoreVectors(items []VectorItem) error {
	for _, item := range items {
		if err := q.StoreVector(item.FileID, item.Vector, item.Description, item.Model); err != nil {
			return fmt.Errorf("批量存储失败，文件ID: %s, 错误: %w", item.FileID, err)
		}
	}
	return nil
}

// SearchSimilar 搜索相似向量
func (q *QdrantClient) SearchSimilar(queryVector []float32, limit int, userID uint, threshold float32, model string) ([]VectorSearchResult, error) {
	return q.SearchVectors(queryVector, limit, userID, threshold)
}

// SearchSimilarWithQuery 带查询的搜索相似向量
func (q *QdrantClient) SearchSimilarWithQuery(queryVector []float32, limit int, userID uint, threshold float32, query string, model string) ([]VectorSearchResult, error) {
	// For now, just delegate to SearchVectors. In a more advanced implementation,
	// we could use the query parameter to filter results
	return q.SearchVectors(queryVector, limit, userID, threshold)
}

func (q *QdrantClient) GetVectorCount(userID uint) (int64, error) {
	return q.GetUserVectorCount(userID)
}

func (q *QdrantClient) GetVector(fileID string) (*models.FileVector, error) {
	// Since we're using direct Qdrant connection, we'd need to implement
	// a way to retrieve the vector from Qdrant and convert it to models.FileVector
	// For now, return an error indicating this is not implemented
	return nil, fmt.Errorf("GetVector not implemented for direct Qdrant connection")
}

// FetchVectorWithPayload 获取指定 fileID 的向量及payload（description/model等）
func (q *QdrantClient) FetchVectorWithPayload(fileID string) ([]float32, map[string]interface{}, error) {
	qdrantID := q.generateQdrantID(fileID)
	// 带上 with_vector 与 with_payload，确保返回完整信息
	url := fmt.Sprintf("%s/collections/%s/points/%s?with_vector=true&with_payload=true", q.baseURL, q.collection, qdrantID)

	resp, err := q.httpClient.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("获取向量失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("获取向量失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var pointResp struct {
		Result struct {
			ID      interface{}            `json:"id"`
			Vector  []float32              `json:"vector"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pointResp); err != nil {
		return nil, nil, fmt.Errorf("解析向量响应失败: %w", err)
	}
	if pointResp.Status != "ok" || len(pointResp.Result.Vector) == 0 {
		return nil, nil, fmt.Errorf("向量不存在或为空")
	}
	return pointResp.Result.Vector, pointResp.Result.Payload, nil
}

func (q *QdrantClient) DeleteVector(fileID string) error {
	// 生成对应的 Qdrant ID
	qdrantID := q.generateQdrantID(fileID)

	reqBody := map[string]interface{}{
		"points": []string{qdrantID},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/collections/%s/points/delete", q.baseURL, q.collection), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除向量请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除向量失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchSimilarByID 通过文件ID搜索相似向量
func (q *QdrantClient) SearchSimilarByID(fileID string, limit int, userID uint, threshold float32, model string) ([]VectorSearchResult, error) {
	// 先从Qdrant获取基准向量
	pointID := q.generateQdrantID(fileID)

	// 构建获取向量的URL，指定需要返回向量数据
	url := fmt.Sprintf("%s/collections/%s/points/%s?with_vector=true", q.baseURL, q.collection, pointID)

	resp, err := q.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取基准向量失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取基准向量失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var pointResp struct {
		Result struct {
			ID      string                 `json:"id"`
			Vector  []float32              `json:"vector"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pointResp); err != nil {
		return nil, fmt.Errorf("解析基准向量响应失败: %w", err)
	}

	if pointResp.Status != "ok" || len(pointResp.Result.Vector) == 0 {
		return nil, fmt.Errorf("基准向量不存在或为空")
	}

	// 使用获取到的向量进行相似度搜索
	return q.SearchVectors(pointResp.Result.Vector, limit, userID, threshold)
}

// IsEnabled 检查是否启用
func (q *QdrantClient) IsEnabled() bool {
	return q.HealthCheck() == nil
}
