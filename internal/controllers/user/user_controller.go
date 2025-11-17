package user

import (
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/activity"
	"pixelpunk/internal/services/user"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	req, err := common.ValidateRequest[dto.RegisterDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID, err := user.RegisterUser(req.Username, req.Email, req.Password, req.Code)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogUserRegister(userID, req.Username, req.Email)

	errors.ResponseSuccess(c, nil, "注册成功")
}

func Login(c *gin.Context) {
	req, err := common.ValidateRequest[dto.LoginDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userInfo, token, err := user.Login(req.Account, req.Password)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	email := ""
	if val, ok := userInfo["email"]; ok {
		if emailStr, isString := val.(string); isString {
			email = emailStr
		}
	}

	data := gin.H{
		"token":    token,
		"userInfo": userInfo,
		"email":    email,
	}

	if userIDVal, ok := userInfo["id"]; ok {
		if userID, isUint := userIDVal.(uint); isUint {
			username := ""
			if usernameVal, ok := userInfo["username"]; ok {
				if usernameStr, isString := usernameVal.(string); isString {
					username = usernameStr
				}
			}

			clientIP := utils.GetClientIP(c)

			activity.LogUserLogin(userID, username, clientIP)
		}
	}

	errors.ResponseSuccess(c, data, "登录成功")
}

func SendRegistrationCode(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SendCodeDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.SendRegistrationCode(req.Email); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}

func SendResetPasswordCode(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SendCodeDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.SendResetPasswordCode(req.Email); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}

func ResetPassword(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ResetPasswordDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "密码重置成功")
}

func UpdatePassword(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ChangePasswordDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := user.UpdatePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogPasswordChange(userID)

	errors.ResponseSuccess(c, nil, "密码修改成功")
}

func UpdateProfile(c *gin.Context) {
	req, err := common.ValidateRequest[dto.UpdateProfileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	userID := middleware.GetCurrentUserID(c)

	updateData := make(map[string]interface{})
	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if c.Request.PostForm.Has("avatar") || c.Request.Header.Get("Content-Type") == "application/json" {
		updateData["avatar"] = req.Avatar
	}
	if c.Request.PostForm.Has("bio") || c.Request.Header.Get("Content-Type") == "application/json" {
		updateData["bio"] = req.Bio
	}
	if c.Request.PostForm.Has("website") || c.Request.Header.Get("Content-Type") == "application/json" {
		updateData["website"] = req.Website
	}
	updatedUser, err := user.UpdateUserProfile(userID, updateData)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	for field := range updateData {
		switch field {
		case "username":
			activity.LogProfileUpdate(userID, "username")
		case "avatar":
			activity.LogProfileUpdate(userID, "avatar")
		case "bio":
			activity.LogProfileUpdate(userID, "bio")
		case "website":
			activity.LogProfileUpdate(userID, "website")
		}
	}

	avatarFullPath := ""
	if updatedUser.Avatar != "" {
		avatarFullPath = utils.GetSystemFileURL(updatedUser.Avatar)
	}

	data := gin.H{
		"id":             updatedUser.ID,
		"username":       updatedUser.Username,
		"email":          updatedUser.Email,
		"avatar":         updatedUser.Avatar,
		"avatarFullPath": avatarFullPath,
		"bio":            updatedUser.Bio,
		"website":        updatedUser.Website,
		"role":           updatedUser.Role,
	}

	errors.ResponseSuccess(c, data, "资料更新成功")
}

func GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	userInfo, err := user.GetUserProfile(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	avatarFullPath := ""
	if userInfo.Avatar != "" {
		avatarFullPath = utils.GetSystemFileURL(userInfo.Avatar)
	} else {
		avatarFullPath = ""
	}

	data := gin.H{
		"id":             userInfo.ID,
		"username":       userInfo.Username,
		"email":          userInfo.Email,
		"avatar":         userInfo.Avatar,
		"avatarFullPath": avatarFullPath,
		"bio":            userInfo.Bio,
		"website":        userInfo.Website,
		"role":           userInfo.Role,
	}

	errors.ResponseSuccess(c, data, "获取成功")
}

func SendChangeEmailCode(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SendChangeEmailCodeDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := user.SendChangeEmailCode(userID, req.NewEmail); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "验证码已发送")
}

func ChangeEmail(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ChangeEmailDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := user.ChangeEmail(userID, req.NewEmail, req.Code); err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogEmailChange(userID, req.NewEmail)

	errors.ResponseSuccess(c, nil, "邮箱更换成功")
}
