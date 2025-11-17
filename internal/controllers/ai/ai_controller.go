package ai

import (
	"pixelpunk/internal/services/ai"
	aiClient "pixelpunk/pkg/ai"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TriggerFileTagging(c *gin.Context) {
	maxFilesStr := c.DefaultQuery("max_files", "0")
	maxFiles, err := strconv.Atoi(maxFilesStr)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的最大文件数量参数"))
		return
	}
	count, err := ai.TriggerFullTaggingProcess(maxFiles)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "启动文件标记处理失败: "+err.Error()))
		return
	}
	errors.ResponseSuccess(c, gin.H{
		"submitted_count": count,
		"status":          "processing",
	}, "文件标记处理已启动")
}

func ResetStuckFiles(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold_minutes", "30")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的阈值时间参数"))
		return
	}
	count, err := ai.ResetStuckPendingFiles(threshold)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "重置卡住的文件失败: "+err.Error()))
		return
	}
	errors.ResponseSuccess(c, gin.H{
		"reset_count": count,
	}, "成功重置卡住的文件")
}

func TestAIConfig(c *gin.Context) {
	var testParams map[string]interface{}
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&testParams); err == nil && len(testParams) > 0 {
			result, err := ai.TestAIConfigurationWithParams(testParams)
			if err != nil {
				errors.HandleError(c, errors.New(errors.CodeInternal, "AI配置测试失败: "+err.Error()))
				return
			}
			if success, ok := result["success"].(bool); ok && success {
				if refreshErr := aiClient.RefreshDefaultClient(); refreshErr != nil {
					errors.HandleError(c, errors.New(errors.CodeInternal, "AI客户端刷新失败: "+refreshErr.Error()))
					return
				}
			}
			errors.ResponseSuccess(c, result, "AI配置测试完成")
			return
		}
	}
	if err := aiClient.RefreshDefaultClient(); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "AI客户端刷新失败: "+err.Error()))
		return
	}
	result, err := ai.TestAIConfiguration()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "AI配置测试失败: "+err.Error()))
		return
	}
	errors.ResponseSuccess(c, result, "AI配置测试完成")
}
