package middleware

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func isStaticAsset(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	staticExts := []string{
		".js", ".css", ".png", ".jpg", ".jpeg", ".gif",
		".svg", ".ico", ".woff", ".woff2", ".ttf", ".eot",
		".pdf", ".mp4", ".webp", ".map",
	}

	for _, staticExt := range staticExts {
		if ext == staticExt {
			return true
		}
	}
	return false
}

func getCacheableExtensions() map[string]int {
	return map[string]int{
		".js":    5184000, // 60天: 60*24*60*60
		".css":   5184000, // 60天
		".png":   5184000, // 60天
		".jpg":   5184000, // 60天
		".jpeg":  5184000, // 60天
		".gif":   5184000, // 60天
		".svg":   5184000, // 60天
		".ico":   5184000, // 60天
		".woff":  5184000, // 60天
		".woff2": 5184000, // 60天
		".ttf":   5184000, // 60天
		".eot":   5184000, // 60天
		".webp":  5184000, // 60天
		".map":   86400,   // 1天: 24*60*60
		".pdf":   2592000, // 30天: 30*24*60*60
		".mp4":   2592000, // 30天
	}
}

func calculateETag(file fs.File) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	modTime := fileInfo.ModTime().Unix()
	size := fileInfo.Size()
	etag := fmt.Sprintf("%x-%x", modTime, size)

	return etag, nil
}

func StaticFileHandler(subFS fs.FS) gin.HandlerFunc {
	cacheableExts := getCacheableExtensions()

	return func(c *gin.Context) {
		reqPath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if reqPath == "" || reqPath == "/" {
			reqPath = "index.html"
		}
		file, err := subFS.Open(reqPath)
		if err != nil {
			if isStaticAsset(reqPath) {
				c.Status(http.StatusNotFound)
				return
			}

			indexFile, indexErr := subFS.Open("index.html")
			if indexErr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			c.Writer.Header().Set("Content-Type", "text/html")
			c.Writer.Header().Set("Cache-Control", "public, max-age=60, must-revalidate")
			c.Writer.Header().Set("Pragma", "no-cache")

			content, readErr := io.ReadAll(indexFile)
			if readErr != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write(content)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		ext := strings.ToLower(filepath.Ext(reqPath))
		switch ext {
		case ".js":
			c.Writer.Header().Set("Content-Type", "application/javascript")
		case ".css":
			c.Writer.Header().Set("Content-Type", "text/css")
		case ".html":
			c.Writer.Header().Set("Content-Type", "text/html")
		case ".svg":
			c.Writer.Header().Set("Content-Type", "image/svg+xml")
		case ".woff":
			c.Writer.Header().Set("Content-Type", "font/woff")
		case ".woff2":
			c.Writer.Header().Set("Content-Type", "font/woff2")
		case ".ttf":
			c.Writer.Header().Set("Content-Type", "font/ttf")
		}

		// HTML 文件特殊处理：使用60秒短缓存，不使用 ETag/Last-Modified 避免浏览器缓存旧版本
		if ext == ".html" {
			c.Writer.Header().Set("Cache-Control", "public, max-age=60, must-revalidate")
			c.Writer.Header().Set("Pragma", "no-cache")
		} else {
			// 其他文件正常处理 ETag 和长期缓存
			etag, err := calculateETag(file)
			if err == nil {
				c.Writer.Header().Set("ETag", fmt.Sprintf(`"%s"`, etag))

				if match := c.GetHeader("If-None-Match"); match != "" {
					if strings.Contains(match, etag) {
						c.Writer.WriteHeader(http.StatusNotModified)
						return
					}
				}
			}

			modTime := fileInfo.ModTime()
			c.Writer.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))

			if ifModSince := c.GetHeader("If-Modified-Since"); ifModSince != "" {
				ifModSinceTime, err := time.Parse(http.TimeFormat, ifModSince)
				if err == nil && modTime.Before(ifModSinceTime.Add(time.Second)) {
					c.Writer.WriteHeader(http.StatusNotModified)
					return
				}
			}

			if maxAge, ok := cacheableExts[ext]; ok {
				c.Writer.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
				c.Writer.Header().Set("Expires", time.Now().Add(time.Duration(maxAge)*time.Second).Format(time.RFC1123))
			} else {
				c.Writer.Header().Set("Cache-Control", "no-cache, must-revalidate")
				c.Writer.Header().Set("Pragma", "no-cache")
			}
		}

		newFile, err := subFS.Open(reqPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		readSeeker, ok := newFile.(io.ReadSeeker)
		if ok {
			modTime := fileInfo.ModTime()
			http.ServeContent(c.Writer, c.Request, reqPath, modTime, readSeeker)
		} else {
			content, err := io.ReadAll(newFile)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))

			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write(content)
		}
	}
}
