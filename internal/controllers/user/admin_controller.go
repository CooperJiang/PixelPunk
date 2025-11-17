package user

import (
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/user"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminGetUserList(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminUserListQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := user.GetUserList(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取用户列表成功")
}

func AdminUpdateUser(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminUpdateUserDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	currentUserID := middleware.GetCurrentUserID(c)

	if err := user.AdminUpdateUser(req, currentUserID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新用户信息成功")
}

func AdminCreateUser(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminCreateUserDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := user.AdminCreateUser(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "创建用户成功")
}

func AdminGetUserDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID不能为空"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID格式不正确"))
		return
	}

	result, err := user.AdminGetUserByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取用户详情成功")
}

func AdminUpdateUserStorage(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminUpdateUserStorageDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	req.CurrentUserID = middleware.GetCurrentUserID(c)

	if err := user.AdminUpdateUserStorage(req); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新存储设置成功")
}

func AdminResetUserPassword(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID不能为空"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID格式不正确"))
		return
	}

	result, err := user.AdminResetUserPassword(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "重置密码成功")
}

func AdminSendUserEmail(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminSendUserEmailDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.AdminSendUserEmail(req); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "邮件发送成功")
}

func AdminToggleUserStatus(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminToggleUserStatusDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.AdminToggleUserStatus(req); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "用户状态更新成功")
}

func AdminDeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID不能为空"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID格式不正确"))
		return
	}

	if err := user.AdminDeleteUser(uint(id)); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除用户成功")
}

func AdminBatchOperateUsers(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminBatchOperateUsersDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.AdminBatchOperateUsers(req); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "批量操作成功")
}
