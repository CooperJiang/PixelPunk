package common

import (
	"fmt"
	"os"
	"sync"
)

// FileDeleteResult 文件删除结果
type FileDeleteResult struct {
	FilePath string // 文件路径
	Success  bool   // 是否成功删除
	Error    error  // 错误信息（如果有）
}

// FileReferenceCheck 是一个函数类型，用于检查文件是否被引用
type FileReferenceCheck func() (bool, int64)

// AsyncDeleteFile 异步删除文件
// filePath: 要删除的文件路径
// isReferenced: 检查文件是否被引用的函数
// referenceCount: 引用计数（如果已知）
// 返回一个通道，可以用来接收删除结果
func AsyncDeleteFile(filePath string, isReferenced FileReferenceCheck, referenceCount int64) <-chan FileDeleteResult {
	resultChan := make(chan FileDeleteResult, 1)

	go func() {
		result := FileDeleteResult{
			FilePath: filePath,
			Success:  false,
		}

		// 如果提供了引用检查函数，则执行检查
		if isReferenced != nil {
			hasReference, count := isReferenced()
			if hasReference {
				// 文件被其他记录引用，跳过删除
				result.Error = fmt.Errorf("文件被引用 (%d 次)", count)
				resultChan <- result
				close(resultChan)
				return
			}
		} else if referenceCount > 0 {
			// 文件被其他记录引用，跳过删除
			result.Error = fmt.Errorf("文件被引用 (%d 次)", referenceCount)
			resultChan <- result
			close(resultChan)
			return
		}

		err := os.Remove(filePath)
		if err != nil {
			if !os.IsNotExist(err) { // 忽略文件不存在的错误
				result.Error = err
			} else {
				// 文件不存在视为成功
				result.Success = true
			}
		} else {
			result.Success = true
		}

		resultChan <- result
		close(resultChan)
	}()

	return resultChan
}

// AsyncDeleteFiles 批量异步删除多个文件
// 返回一个通道，可以用来接收所有删除结果
func AsyncDeleteFiles(filePaths []string, checkReferenceFn FileReferenceCheck) <-chan []FileDeleteResult {
	resultChan := make(chan []FileDeleteResult, 1)

	go func() {
		var results []FileDeleteResult
		var wg sync.WaitGroup
		resultsMutex := &sync.Mutex{}

		for _, path := range filePaths {
			if path == "" {
				continue
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()

				resCh := AsyncDeleteFile(filePath, checkReferenceFn, 0)
				res := <-resCh

				// 线程安全地添加结果
				resultsMutex.Lock()
				results = append(results, res)
				resultsMutex.Unlock()
			}(path)
		}

		wg.Wait()
		resultChan <- results
		close(resultChan)
	}()

	return resultChan
}

// filePath: 要删除的文件路径
// referenceCount: 已知的引用计数
// 如果文件被引用或删除失败，返回错误
func DeleteFileWithCheck(filePath string, referenceCount int64) error {
	if referenceCount > 0 {
		return fmt.Errorf("文件 %s 被其他 %d 个记录引用，跳过删除", filePath, referenceCount)
	}

	if filePath == "" {
		return nil // 空路径，不需要操作
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			// 文件不存在视为成功
			return nil
		}
		return fmt.Errorf("删除文件失败 %s: %w", filePath, err)
	}

	return nil
}
