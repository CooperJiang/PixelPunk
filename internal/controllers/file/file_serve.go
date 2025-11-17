package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/utils"

	"github.com/gin-gonic/gin"
)

func ServeFile(c *gin.Context) {
	// 从 context 中获取中间件设置的文件信息
	fileObj, exists := c.Get("file_info")
	if !exists {
		// 文件不存在，返回默认文件
		assets.ServeDefaultFile(c, assets.FileTypeNotFound)
		return
	}

	// 类型断言，将 interface{} 转换为 models.File
	fileInfo, ok := fileObj.(models.File)
	if !ok {
		// 文件信息无效，返回默认文件
		assets.ServeDefaultFile(c, assets.FileTypeNotFound)
		return
	}

	forceThumbnail := false
	if value, exists := c.Get("forceThumbnail"); exists {
		forceThumbnail, _ = value.(bool)
	}

	// 根据存储类型处理文件访问
	result, isLocal, isProxy, err := filesvc.ServeFile(fileInfo, forceThumbnail)
	if err != nil {
		if customErr, ok := err.(*errors.Error); ok && customErr.Code == errors.CodeFileNotFound {
			// 文件不存在，返回默认文件
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}
		// 其他错误仍然按原方式处理
		errors.HandleError(c, err)
		return
	}

	// 根据结果类型决定响应方式
	if isLocal {
		localPath := result.(string)
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			// 本地文件不存在，返回默认文件
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}
		c.File(localPath)
	} else if isProxy {
		proxyResp := result.(*filesvc.ProxyResponse)
		defer proxyResp.Content.Close()

		// 设置Content-Type
		c.Header("Content-Type", proxyResp.ContentType)
		// 设置Content-Length以支持真实下载进度
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

		c.Status(http.StatusOK)
		io.Copy(c.Writer, proxyResp.Content)
	} else {
		url := result.(string)
		c.Redirect(http.StatusFound, url)
	}
}

func ServeThumbnailFile(c *gin.Context) {
	// 从 context 中获取中间件设置的文件信息
	fileObj, exists := c.Get("file_info")
	if !exists {
		// 缩略图不存在，返回默认文件
		assets.ServeDefaultFile(c, assets.FileTypeNotFound)
		return
	}

	// 类型断言，将 interface{} 转换为 models.File
	fileInfo, ok := fileObj.(models.File)
	if !ok {
		// 缩略图信息无效，返回默认文件
		assets.ServeDefaultFile(c, assets.FileTypeNotFound)
		return
	}

	// 根据存储类型处理文件访问
	result, isLocal, isProxy, err := filesvc.ServeFile(fileInfo, true)
	if err != nil {
		if customErr, ok := err.(*errors.Error); ok && customErr.Code == errors.CodeFileNotFound {
			// 缩略图不存在，返回默认文件
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}
		// 其他错误仍然按原方式处理
		errors.HandleError(c, err)
		return
	}

	// 根据结果类型决定响应方式
	if isLocal {
		localPath := result.(string)
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			// 本地文件不存在，返回默认文件
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}
		c.File(localPath)
	} else if isProxy {
		proxyResp := result.(*filesvc.ProxyResponse)
		defer proxyResp.Content.Close()

		// 设置Content-Type
		c.Header("Content-Type", proxyResp.ContentType)
		// 对于缩略图，如果有实际内容长度则设置，否则不设置避免大小不匹配
		if proxyResp.ContentLength > 0 {
			c.Header("Content-Length", strconv.FormatInt(proxyResp.ContentLength, 10))
		}

		c.Status(http.StatusOK)
		io.Copy(c.Writer, proxyResp.Content)
	} else {
		url := result.(string)
		c.Redirect(http.StatusFound, url)
	}
}

func GenerateFileLink(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}
	expireMinutes, _ := strconv.Atoi(c.DefaultQuery("expire", "5"))
	if expireMinutes <= 0 || expireMinutes > 60 {
		expireMinutes = 5
	}
	res, err := filesvc.GenerateTemporaryLink(userID, fileID, expireMinutes)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, gin.H{
		"file_id":       res.FileID,
		"temporary_url": res.URL,
		"expires_in":    res.ExpiresIn,
		"expires_at":    res.ExpiresAt.Format(time.RFC3339),
	}, "生成临时链接成功")
}

func GetFileShare(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "文件不存在或无权限访问"))
		return
	}
}

func DownloadFile(c *gin.Context) {
	// 从中间件中获取文件信息（已经过权限验证）
	fileObj, exists := c.Get("file_info")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件信息获取失败"))
		return
	}

	file, ok := fileObj.(models.File)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件信息无效"))
		return
	}

	fileID := file.ID

	// 获取质量参数，默认为原图
	quality := c.DefaultQuery("quality", "original")
	isThumb := quality != "original"

	// 获取当前用户（对于公开文件可能为nil）
	var currentUserID uint
	if currentUser := middleware.GetCurrentUser(c); currentUser != nil {
		currentUserID = currentUser.UserID
	}

	// 设置下载响应头的基础文件名（先取显示名，空则取原名）
	fileName := file.DisplayName
	if fileName == "" {
		fileName = file.OriginalName
	}

	// 确保文件名安全（移除危险字符）
	fileName = utils.GetSafeFilename(fileName)

	// 确保文件名有正确的扩展名
	if fileName != "" {
		ext := filepath.Ext(fileName)
		if ext == "" {
			// 如果没有扩展名，使用format字段添加扩展名
			if file.Format != "" {
				fileName = fileName + "." + filesvc.GetCorrectFileExtension(file.Format)
			}
		} else {
			// 如果有扩展名但不正确，替换为正确的扩展名
			correctExt := filesvc.GetCorrectFileExtension(file.Format)
			if correctExt != "" && !strings.EqualFold(ext, "."+correctExt) {
				nameWithoutExt := strings.TrimSuffix(fileName, ext)
				fileName = nameWithoutExt + "." + correctExt
			}
		}
	}

	// 根据quality参数获取相应的文件文件
	result, isLocal, isProxy, err := filesvc.ServeFile(file, isThumb)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	go func() {
		downloadLog := &models.FileDownloadLog{
			UserID:    currentUserID, // 用户ID（公开文件下载时可能为0）
			FileID:    fileID,
			FileSize:  file.Size,
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}
		if err := database.DB.Create(downloadLog).Error; err != nil {
			logger.Error("记录下载日志失败: %v", err)
		}
	}()

	// 根据quality参数调整文件名
	if isThumb && quality != "" && quality != "original" {
		// 为缩略图或其他质量版本添加后缀
		ext := filepath.Ext(fileName)
		name := strings.TrimSuffix(fileName, ext)
		fileName = fmt.Sprintf("%s_%s%s", name, quality, ext)
	}

	// 设置正确的Content-Disposition头，支持中文文件名
	c.Header("Content-Disposition", utils.SetContentDispositionFilename(fileName))
	c.Header("Accept-Ranges", "bytes")

	switch {
	case isLocal:
		filePath := result.(string)

		// 检查文件是否为ASCII格式，如果是则转换后返回
		if fileData, err := os.ReadFile(filePath); err == nil {
			// 检查是否为真正的ASCII数组格式 (如 "[60, 63, 120, ...]")
			isAsciiArray := false
			if len(fileData) > 10 && fileData[0] == '[' {
				// 进一步检查是否包含数字和逗号的模式
				str := string(fileData[:100]) // 检查前100个字符
				if strings.Contains(str, ", ") && (strings.Contains(str, "60") || strings.Contains(str, "255")) {
					isAsciiArray = true
				}
			}

			if isAsciiArray {
				normalizedData := adapter.NormalizePossiblyTextualBytes(fileData, "[DOWNLOAD]")
				if len(normalizedData) != len(fileData) {
					// 根据格式设置正确的Content-Type（保持DB格式）
					contentType := filesvc.GetContentTypeByFormat(file.Format)
					c.Header("Content-Type", contentType)
					c.Data(http.StatusOK, contentType, normalizedData)
					return
				}
			}

			// 设置正确的Content-Type（保持DB格式）
			contentType := filesvc.GetContentTypeByFormat(file.Format)
			c.Header("Content-Type", contentType)
		}

		c.File(filePath)
	case isProxy:
		proxyResp := result.(*filesvc.ProxyResponse)
		defer proxyResp.Content.Close()
		// 设置Content-Length以支持真实下载进度
		c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
		c.Status(http.StatusOK)
		io.Copy(c.Writer, proxyResp.Content)
	default:
		c.Redirect(http.StatusTemporaryRedirect, result.(string))
	}
}

// helpers moved to services/file (GetCorrectFileExtension, GetContentTypeByFormat)
