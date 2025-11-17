package vector

import (
	"pixelpunk/internal/models"
	"time"
)

// VectorStorage 向量存储接口
type VectorStorage interface {
	StoreVector(fileID string, vector []float32, description string, model string) error
	GetVector(fileID string) (*models.FileVector, error)
	SearchSimilar(queryVector []float32, limit int, userID uint, threshold float32, model string) ([]VectorSearchResult, error)
	SearchSimilarWithQuery(queryVector []float32, limit int, userID uint, threshold float32, query string, model string) ([]VectorSearchResult, error)
	DeleteVector(fileID string) error
	BatchStoreVectors(items []VectorItem) error
	GetVectorCount(userID uint) (int64, error)
	GetStorageStats() (*VectorStorageStats, error)
	VectorExists(fileID string) (bool, error) // 新增：检查向量是否存在
}

// VectorItem 批量向量处理项
type VectorItem struct {
	FileID      string
	Vector      []float32
	Description string
	Model       string
}

// VectorSearchResult 向量搜索结果
type VectorSearchResult struct {
	FileID      string  `json:"file_id"`
	Similarity  float32 `json:"similarity"`
	Description string  `json:"description"`
	Score       float32 `json:"score"` // 归一化分数 (0-100)
}

// VectorStorageStats 向量存储统计信息
type VectorStorageStats struct {
	TotalVectors   int64     `json:"total_vectors"`
	CompletedCount int64     `json:"completed_count"`
	PendingCount   int64     `json:"pending_count"`
	FailedCount    int64     `json:"failed_count"`
	AverageNorm    float64   `json:"average_norm"`
	LastUpdateTime time.Time `json:"last_update_time"`
	StorageSize    int64     `json:"storage_size_bytes"`
}
