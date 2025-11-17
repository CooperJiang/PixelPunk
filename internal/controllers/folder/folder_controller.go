package folder

import (
	"pixelpunk/internal/controllers/folder/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/activity"
	"pixelpunk/internal/services/folder"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func CreateFolder(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.CreateFolderDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	folderInfo, err := folder.CreateFolder(userID, req.Name, req.ParentID, req.Permission, req.Description)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogFolderCreate(userID, req.Name)

	errors.ResponseSuccess(c, folderInfo, "创建成功")
}

func ListFolders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	parentID := c.Query("parent_id")

	folders, err := folder.ListFolders(userID, parentID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, folders, "获取成功")
}

func GetFolderDetail(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	folderID := c.Param("folder_id")
	if folderID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空"))
		return
	}

	folderDetail, err := folder.GetFolderDetail(userID, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, folderDetail, "获取成功")
}

func UpdateFolder(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.UpdateFolderDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	oldFolderInfo, err := folder.GetFolderDetail(userID, req.FolderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	folderInfo, err := folder.UpdateFolder(userID, req.FolderID, req.Name, req.ParentID, req.Permission, req.Description)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.Name != "" && req.Name != oldFolderInfo.Name {
		activity.LogFolderRename(userID, oldFolderInfo.Name, req.Name)
	}

	errors.ResponseSuccess(c, folderInfo, "更新成功")
}

func DeleteFolder(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	folderID := c.Param("folder_id")
	if folderID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空"))
		return
	}

	folderInfo, err := folder.GetFolderDetail(userID, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = folder.DeleteFolder(userID, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogFolderDelete(userID, folderInfo.Name, int(folderInfo.FileCount))

	errors.ResponseSuccess(c, gin.H{"id": folderID}, "删除成功")
}

func ListFolderContents(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.FolderContentsDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.SortBy == "" {
		req.SortBy = "custom"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	contents, err := folder.ListFolderContentsWithoutPagination(
		userID,
		req.ParentID,
		req.Keyword,
		req.AccessLevel,
		req.SortBy,
		req.SortOrder,
	)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, contents, "获取成功")
}

func ToggleAccessLevel(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	folderID := c.Param("folder_id")
	if folderID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空"))
		return
	}

	oldFolderInfo, err := folder.GetFolderDetail(userID, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	folderInfo, err := folder.ToggleFolderAccessLevel(userID, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogFolderAccessLevelChange(userID, folderInfo.Name, oldFolderInfo.Permission, folderInfo.Permission)

	errors.ResponseSuccess(c, folderInfo, "权限切换成功")
}

func ReorderFolders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.ReorderFoldersRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = folder.ReorderFolders(userID, req.ParentID, req.FolderIDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "文件夹排序更新成功")
}

func GetFolderTree(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	tree, err := folder.GetFolderTree(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, tree, "获取文件夹树成功")
}

func SearchFolders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyword := c.Query("keyword")
	if keyword == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "搜索关键词不能为空"))
		return
	}

	folders, err := folder.SearchFolders(userID, keyword)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, folders, "搜索成功")
}

func GetFolderPathChain(c *gin.Context) {
	folderID := c.Param("folder_id")
	if folderID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空"))
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	pathChain, err := folder.GetFolderPathChain(folderID, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, pathChain, "获取文件夹路径链成功")
}

func GetBatchFolderPathChains(c *gin.Context) {
	req, err := common.ValidateRequest[dto.BatchFolderPathChainsRequestDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	if len(req.FolderIDs) > 50 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "批量查询文件夹数量不能超过50个"))
		return
	}

	pathChains, err := folder.GetBatchFolderPathChains(req.FolderIDs, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, pathChains, "批量获取文件夹路径链成功")
}
