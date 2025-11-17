package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
)

func GetFolderTree(userID uint) ([]*TreeNodeResponse, error) {
	var folders []models.Folder
	if err := database.DB.Where("user_id = ?", userID).Order("sort_order ASC, name ASC").Find(&folders).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹列表失败")
	}
	return buildFolderTree(folders), nil
}

func buildFolderTree(folders []models.Folder) []*TreeNodeResponse {
	nodeMap := make(map[string]*TreeNodeResponse)
	var rootNodes []*TreeNodeResponse
	for _, folder := range folders {
		folderResponse := toResponse(&folder)
		node := &TreeNodeResponse{ID: folder.ID, Label: folder.Name, Icon: "fa-folder", Count: 0, Level: folderResponse.Level, Children: make([]*TreeNodeResponse, 0), Data: folderResponse}
		nodeMap[folder.ID] = node
	}
	for _, folder := range folders {
		node := nodeMap[folder.ID]
		if folder.ParentID == "" || folder.ParentID == "0" {
			rootNodes = append(rootNodes, node)
		} else if parentNode, exists := nodeMap[folder.ParentID]; exists {
			parentNode.Children = append(parentNode.Children, node)
		}
	}
	for _, node := range nodeMap {
		var imageCount int64
		database.DB.Model(&models.File{}).Where("folder_id = ?", node.ID).Where("status <> ?", "pending_deletion").Count(&imageCount)
		node.Count = imageCount
	}
	return rootNodes
}
