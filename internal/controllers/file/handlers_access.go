package file

import (
	"io"
	"pixelpunk/internal/models"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ServeFileByID(c *gin.Context) {
	fileObj, exists := c.Get("file_info")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件信息未找到"))
		return
	}

	fileInfo, ok := fileObj.(models.File)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeInternal, "文件信息类型转换失败"))
		return
	}

	serveFileByInfo(c, fileInfo, false)
}

func ServeThumbByID(c *gin.Context) {
	fileObj, exists := c.Get("file_info")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件信息未找到"))
		return
	}

	fileInfo, ok := fileObj.(models.File)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeInternal, "文件信息类型转换失败"))
		return
	}

	serveFileByInfo(c, fileInfo, true)
}

func ServeFileByShortURL(c *gin.Context) {
	fileObj, exists := c.Get("file_info")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件信息未找到"))
		return
	}

	fileInfo, ok := fileObj.(models.File)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeInternal, "文件信息类型转换失败"))
		return
	}

	serveFileByInfo(c, fileInfo, false)
}

func serveFileByInfo(c *gin.Context, fileInfo models.File, isThumb bool) {
	result, isLocalPath, isProxy, err := filesvc.ServeFile(fileInfo, isThumb)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Header("Cache-Control", "public, max-age=2592000, immutable")
	c.Header("Access-Control-Allow-Origin", "*")

	if isLocalPath {
		if filePath, ok := result.(string); ok {
			c.File(filePath)
		}
		return
	}

	if isProxy {
		if proxyResp, ok := result.(*filesvc.ProxyResponse); ok {
			defer proxyResp.Content.Close()

			c.Header("Content-Type", proxyResp.ContentType)

			if isThumb {
				if proxyResp.ContentLength > 0 {
					c.Header("Content-Length", strconv.FormatInt(proxyResp.ContentLength, 10))
				}
			} else {
				c.Header("Content-Length", strconv.FormatInt(fileInfo.Size, 10))
			}

			c.Status(200)
			io.Copy(c.Writer, proxyResp.Content)
		}
		return
	}

	if url, ok := result.(string); ok {
		c.Redirect(302, url)
	}
}
