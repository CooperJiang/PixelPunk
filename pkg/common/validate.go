package common

import (
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string
	Message string
}

func ValidateRequest[T any](c *gin.Context) (*T, error) {
	var req T
	var messages map[string]string
	if v, ok := any(&req).(interface{ GetValidationMessages() map[string]string }); ok {
		messages = v.GetValidationMessages()
	}
	var err error
	method := c.Request.Method
	contentType := c.ContentType()
	if method == "GET" {
		err = c.ShouldBindQuery(&req)
	} else if contentType == "application/json" {
		err = c.ShouldBindJSON(&req)
	} else if contentType == "multipart/form-data" {
		err = c.ShouldBindWith(&req, binding.FormMultipart)
	} else {
		err = c.ShouldBind(&req)
	}
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				field := validationErr.StructField()
				tag := validationErr.Tag()
				msgKey := field + "." + tag
				if msg, exists := messages[msgKey]; exists {
					return nil, errors.NewValidationError(field, msg)
				}
			}
			return nil, errors.New(errors.CodeInvalidParameter, "请求参数验证失败")
		}
		return nil, errors.New(errors.CodeInvalidParameter, "无效的请求参数")
	}
	return &req, nil
}
