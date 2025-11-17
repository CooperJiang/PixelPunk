package errors

import (
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternalServer = &Error{
		Code:    CodeInternal,
		Message: "服务器内部错误",
		Detail:  "服务器内部错误",
		Time:    time.Now(),
	}

	ErrInvalidParams = &Error{
		Code:    CodeInvalidParameter,
		Message: "无效的请求参数",
		Detail:  "无效的请求参数",
		Time:    time.Now(),
	}

	ErrUnauthorized = &Error{
		Code:    CodeUnauthorized,
		Message: "未授权访问",
		Detail:  "未授权访问",
		Time:    time.Now(),
	}

	ErrForbidden = &Error{
		Code:    CodeForbidden,
		Message: "禁止访问",
		Detail:  "禁止访问",
		Time:    time.Now(),
	}

	ErrNotFound = &Error{
		Code:    CodeNotFound,
		Message: "资源不存在",
		Detail:  "资源不存在",
		Time:    time.Now(),
	}
)

func ResponseError(c *gin.Context, err *Error, detail string) {
	if detail != "" {
		err.Detail = detail
	}
	requestID, exists := c.Get("RequestID")
	if exists {
		err.RequestID = requestID.(string)
	}
	response := Response{
		Code:      int(err.Code),
		Message:   err.Message,
		RequestID: err.RequestID,
		Timestamp: time.Now().Unix(),
	}
	statusCode := HTTPStatus(err)
	c.JSON(statusCode, response)
}

func ResponseErrorWithMessage(c *gin.Context, err *Error, message string) {
	customErr := *err
	if message != "" {
		customErr.Message = message
	}
	requestID, exists := c.Get("RequestID")
	if exists {
		customErr.RequestID = requestID.(string)
	}
	response := Response{
		Code:      int(customErr.Code),
		Message:   customErr.Message,
		RequestID: customErr.RequestID,
		Timestamp: time.Now().Unix(),
	}
	statusCode := HTTPStatus(err)
	c.JSON(statusCode, response)
}
