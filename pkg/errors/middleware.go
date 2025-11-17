package errors

import (
	"fmt"
	"net/http"
	"pixelpunk/pkg/logger"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				logger.Error("[PANIC] %v\nTrace: %s", r, stackTrace)
				err := &Error{
					Code:      CodeInternal,
					Message:   "服务器内部错误",
					Detail:    fmt.Sprintf("%v", r),
					Stack:     stackTrace,
					Time:      time.Now(),
					RequestID: requestID,
				}
				responseError(c, err)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func responseError(c *gin.Context, err error) {
	var apiErr *Error
	var statusCode int
	if e, ok := err.(*Error); ok {
		apiErr = e
		statusCode = HTTPStatus(e)
	} else {
		apiErr = &Error{
			Code:    CodeInternal,
			Message: "服务器内部错误",
			Detail:  err.Error(),
			Time:    time.Now(),
		}
		statusCode = http.StatusInternalServerError
	}
	requestID, exists := c.Get("RequestID")
	if exists {
		apiErr.RequestID = requestID.(string)
	}
	response := Response{
		Code:      int(apiErr.Code),
		Message:   apiErr.Message,
		RequestID: apiErr.RequestID,
		Timestamp: time.Now().Unix(),
	}
	c.JSON(statusCode, response)
}

func ResponseSuccess(c *gin.Context, data interface{}, message string) {
	requestID, _ := c.Get("RequestID")
	response := Response{
		Code:      200,
		Message:   message,
		Data:      data,
		RequestID: requestID.(string),
		Timestamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, response)
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	responseError(c, err)
}
