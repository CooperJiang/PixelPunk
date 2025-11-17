package file

import (
	"os"
	"path/filepath"
	"strconv"

	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "头像上传失败: "+err.Error()))
		return
	}

	publicURL, _, err := filesvc.UploadAvatar(file)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"avatar_url": publicURL,
		"full_url":   utils.GetSystemFileURL(publicURL), // 使用工具函数获取完整URL
	}, "头像上传成功")
}

func ServeAvatarFile(c *gin.Context) {
	fileName := c.Param("fileName")
	if fileName == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件名不能为空"))
		return
	}

	avatarPath := filepath.Join(filesvc.AvatarUploadDir, fileName)
	if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "头像文件不存在"))
		return
	}

	c.File(avatarPath)
}

func UploadAdminFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件上传失败: "+err.Error()))
		return
	}

	relativePath, fullURL, err := filesvc.UploadAdminFile(file)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"file_url": relativePath,
		"full_url": fullURL,
	}, "文件上传成功")
}

func ServeAdminFile(c *gin.Context) {
	fileName := c.Param("fileName")
	if fileName == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件名不能为空"))
		return
	}

	filePath := filepath.Join(filesvc.FileUploadDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件不存在"))
		return
	}

	c.File(filePath)
}

func UserGetTagList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	tags, total, err := filesvc.AdminGetTagList(page, size, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": tags,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回标签列表
		data = gin.H{
			"items": tags,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取标签列表成功")
}

func UserGetColorList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	colors, total, err := filesvc.AdminGetColorList(page, size, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": colors,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回颜色列表
		data = gin.H{
			"items": colors,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取颜色列表成功")
}

func GuestGetTagList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	// 调用服务层获取标签列表，传入特殊参数表示只查询推荐文件
	// 第一个参数0表示不限制用户，第二个参数1表示只查询推荐文件
	tags, total, err := filesvc.AdminGetTagList(page, size, 0, 1)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": tags,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回标签列表
		data = gin.H{
			"items": tags,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取推荐标签列表成功")
}

func GetPublicTagList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	keyword := c.Query("keyword") // 可选的关键词搜索

	tags, total, err := filesvc.GetPublicTagList(page, size, keyword)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": tags,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回标签列表
		data = gin.H{
			"items": tags,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取公开标签列表成功")
}

func GuestGetColorList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	// 调用服务层获取颜色列表，传入特殊参数表示只查询推荐文件
	// 第一个参数0表示不限制用户，第二个参数1表示只查询推荐文件
	colors, total, err := filesvc.AdminGetColorList(page, size, 0, 1)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": colors,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回颜色列表
		data = gin.H{
			"items": colors,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取推荐颜色列表成功")
}

func UserGetCategoryList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	db := database.GetDB()

	var categories []models.FileCategory
	if err := db.Model(&models.FileCategory{}).Where("user_id = ?", userID).Order("sort_order ASC").Find(&categories).Error; err != nil {
		logger.Error("获取分类列表失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "获取分类列表失败"))
		return
	}

	items := make([]gin.H, len(categories))
	for i, cat := range categories {
		items[i] = gin.H{
			"id":         cat.ID,
			"name":       cat.Name,
			"file_count": cat.FileCount,
		}
	}

	errors.ResponseSuccess(c, gin.H{
		"items": items,
		"total": len(categories),
	}, "获取分类列表成功")
}
