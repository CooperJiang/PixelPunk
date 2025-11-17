package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

// HashResult 哈希计算结果
type HashResult struct {
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256,omitempty"`
	Size   int64  `json:"size"`
}

// CalculateFileMD5 计算文件MD5哈希值
func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// CalculateReaderMD5 计算Reader的MD5哈希值
func CalculateReaderMD5(reader io.Reader) (string, error) {
	hasher := md5.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("读取数据失败: %w", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// CalculateDataMD5 计算字节数据的MD5哈希值
func CalculateDataMD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// CalculateFileHashes 计算文件的多种哈希值
func CalculateFileHashes(filePath string) (*HashResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return CalculateReaderHashes(file, stat.Size())
}

// CalculateReaderHashes 计算Reader的多种哈希值
func CalculateReaderHashes(reader io.Reader, size int64) (*HashResult, error) {
	md5Hasher := md5.New()
	sha256Hasher := sha256.New()

	// 使用MultiWriter同时计算两种哈希
	multiWriter := io.MultiWriter(md5Hasher, sha256Hasher)

	written, err := io.Copy(multiWriter, reader)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}

	// 如果size为0，使用实际读取的大小
	if size == 0 {
		size = written
	}

	return &HashResult{
		MD5:    fmt.Sprintf("%x", md5Hasher.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256Hasher.Sum(nil)),
		Size:   size,
	}, nil
}

// CalculateDataHashes 计算字节数据的多种哈希值
func CalculateDataHashes(data []byte) *HashResult {
	md5Hasher := md5.New()
	sha256Hasher := sha256.New()

	md5Hasher.Write(data)
	sha256Hasher.Write(data)

	return &HashResult{
		MD5:    fmt.Sprintf("%x", md5Hasher.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256Hasher.Sum(nil)),
		Size:   int64(len(data)),
	}
}

// VerifyFileHash 验证文件哈希值
func VerifyFileHash(filePath, expectedHash string, hashType string) (bool, error) {
	var actualHash string
	var err error

	switch hashType {
	case "md5":
		actualHash, err = CalculateFileMD5(filePath)
	case "sha256":
		result, hashErr := CalculateFileHashes(filePath)
		if hashErr != nil {
			err = hashErr
		} else {
			actualHash = result.SHA256
		}
	default:
		return false, fmt.Errorf("不支持的哈希类型: %s", hashType)
	}

	if err != nil {
		return false, err
	}

	return actualHash == expectedHash, nil
}

// StreamHasher 流式哈希计算器
type StreamHasher struct {
	hashers []hash.Hash
	types   []string
}

func NewStreamHasher(hashTypes ...string) *StreamHasher {
	sh := &StreamHasher{
		hashers: make([]hash.Hash, 0, len(hashTypes)),
		types:   make([]string, 0, len(hashTypes)),
	}

	for _, hashType := range hashTypes {
		switch hashType {
		case "md5":
			sh.hashers = append(sh.hashers, md5.New())
			sh.types = append(sh.types, "md5")
		case "sha256":
			sh.hashers = append(sh.hashers, sha256.New())
			sh.types = append(sh.types, "sha256")
		}
	}

	return sh
}

// Write 实现io.Writer接口
func (sh *StreamHasher) Write(p []byte) (n int, err error) {
	for _, hasher := range sh.hashers {
		hasher.Write(p)
	}
	return len(p), nil
}

// Sum 获取所有哈希值
func (sh *StreamHasher) Sum() map[string]string {
	result := make(map[string]string)
	for i, hasher := range sh.hashers {
		result[sh.types[i]] = fmt.Sprintf("%x", hasher.Sum(nil))
	}
	return result
}

// Reset 重置所有哈希计算器
func (sh *StreamHasher) Reset() {
	for _, hasher := range sh.hashers {
		hasher.Reset()
	}
}
