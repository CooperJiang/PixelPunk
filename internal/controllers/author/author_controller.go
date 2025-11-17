package author

import (
	"pixelpunk/internal/services/author"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAuthorHomepage(c *gin.Context) {
	authorIDStr := c.Param("author_id")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的作者ID"))
		return
	}

	homepage, err := author.GetAuthorHomepage(uint(authorID))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, homepage, "获取作者主页信息成功")
}

func GetAuthorFolder(c *gin.Context) {
	authorIDStr := c.Param("author_id")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的作者ID"))
		return
	}

	folderID := c.Param("folder_id")
	if folderID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	folderData, err := author.GetAuthorFolder(uint(authorID), folderID, page, size)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, folderData, "获取文件夹内容成功")
}
