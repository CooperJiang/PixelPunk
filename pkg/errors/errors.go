package errors

import (
	"fmt"
	"pixelpunk/pkg/logger"
	"runtime"
	"strings"
	"time"
)

type ErrorCode int

const (
	CodeUnknown            ErrorCode = 1
	CodeInternal           ErrorCode = 2
	CodeInvalidParameter   ErrorCode = 100
	CodeInvalidRequest     ErrorCode = 101
	CodeUnauthorized       ErrorCode = 102
	CodeForbidden          ErrorCode = 103
	CodeNotFound           ErrorCode = 104
	CodeMethodNotAllowed   ErrorCode = 105
	CodeTimeout            ErrorCode = 106
	CodeConflict           ErrorCode = 107
	CodeRateLimited        ErrorCode = 108
	CodeValidationFailed   ErrorCode = 109
	CodeServiceUnavailable ErrorCode = 110

	CodeUserNotFound      ErrorCode = 1000
	CodeWrongPassword     ErrorCode = 1001
	CodeUserDisabled      ErrorCode = 1002
	CodeUserExists        ErrorCode = 1003
	CodeInvalidAuthToken  ErrorCode = 1004
	CodeExpiredAuthToken  ErrorCode = 1005
	CodeInvalidVerifyCode ErrorCode = 1006
	CodeEmailExists       ErrorCode = 1007
	CodeEmailSendFailed   ErrorCode = 1008
	CodeInvalidToken      ErrorCode = 1009
	CodeTokenExpired      ErrorCode = 1010
	CodeTokenUsed         ErrorCode = 1011

	CodeDBConnectionFailed ErrorCode = 2000
	CodeDBQueryFailed      ErrorCode = 2001
	CodeDBCreateFailed     ErrorCode = 2002
	CodeDBUpdateFailed     ErrorCode = 2003
	CodeDBDeleteFailed     ErrorCode = 2004
	CodeDBNoRecord         ErrorCode = 2005
	CodeDBDuplicate        ErrorCode = 2006
	CodeDBTransaction      ErrorCode = 2007
	CodeDBCommitFailed     ErrorCode = 2008

	CodeThirdPartyService ErrorCode = 3000
	CodeRedisError        ErrorCode = 3001
	CodeEmailServiceError ErrorCode = 3002
	CodeSMSServiceError   ErrorCode = 3003
	CodePaymentError      ErrorCode = 3004

	CodeFileTooLarge            ErrorCode = 4000
	CodeFileTypeNotSupported    ErrorCode = 4001
	CodeFileUploadFailed        ErrorCode = 4002
	CodeFileNotFound            ErrorCode = 4003
	CodeFileDeleteFailed        ErrorCode = 4004
	CodeFileUpdateFailed        ErrorCode = 4005
	CodeFileDownloadFailed      ErrorCode = 4006
	CodeFileAccessDenied        ErrorCode = 4007
	CodeStorageLimitExceeded    ErrorCode = 4008
	CodeBandwidthLimitExceeded  ErrorCode = 4009
	CodeUploadLimitExceeded     ErrorCode = 4010
	CodeFileNotOwned            ErrorCode = 4011
	CodeFileFormatNotSupport    ErrorCode = 4012
	CodeStorageProviderNotFound ErrorCode = 4013
	CodeUploadSessionExpired    ErrorCode = 4014
	CodeChunkUploadFailed       ErrorCode = 4015
	CodeChunkMergeError         ErrorCode = 4016

	CodeFolderNotFound      ErrorCode = 5000
	CodeFolderCreateFailed  ErrorCode = 5001
	CodeFolderDeleteFailed  ErrorCode = 5002
	CodeFolderUpdateFailed  ErrorCode = 5003
	CodeFolderNotEmpty      ErrorCode = 5004
	CodeFolderNameDuplicate ErrorCode = 5005
	CodeInvalidFolderName   ErrorCode = 5006

	CodeIPAccessDenied       ErrorCode = 6000
	CodeIPNotInWhitelist     ErrorCode = 6001
	CodeIPInBlacklist        ErrorCode = 6002
	CodeRefererAccessDenied  ErrorCode = 6003
	CodeDomainNotInWhitelist ErrorCode = 6004
	CodeDomainInBlacklist    ErrorCode = 6005

	CodeSystemNotInstalled ErrorCode = 7001
	CodeInstallationFailed ErrorCode = 7002
	CodeConfigurationError ErrorCode = 7003

	CodeShareExpired        ErrorCode = 8001
	CodeShareAccessExceeded ErrorCode = 8002
	CodeShareNotFound       ErrorCode = 8003
	CodeSharePasswordWrong  ErrorCode = 8004
)

var errorCodeToHTTPStatus = map[ErrorCode]int{
	CodeUnknown:            500,
	CodeInternal:           500,
	CodeInvalidParameter:   400,
	CodeInvalidRequest:     400,
	CodeUnauthorized:       401,
	CodeForbidden:          403,
	CodeNotFound:           404,
	CodeMethodNotAllowed:   405,
	CodeTimeout:            408,
	CodeConflict:           409,
	CodeRateLimited:        429,
	CodeValidationFailed:   400,
	CodeServiceUnavailable: 503,

	CodeUserNotFound:      400,
	CodeWrongPassword:     400,
	CodeUserDisabled:      403,
	CodeUserExists:        409,
	CodeInvalidAuthToken:  401,
	CodeExpiredAuthToken:  401,
	CodeInvalidVerifyCode: 400,
	CodeEmailExists:       409,
	CodeEmailSendFailed:   500,
	CodeInvalidToken:      400,
	CodeTokenExpired:      400,
	CodeTokenUsed:         400,

	CodeDBConnectionFailed: 500,
	CodeDBQueryFailed:      500,
	CodeDBCreateFailed:     500,
	CodeDBUpdateFailed:     500,
	CodeDBDeleteFailed:     500,
	CodeDBNoRecord:         404,
	CodeDBDuplicate:        409,
	CodeDBTransaction:      500,
	CodeDBCommitFailed:     500,

	CodeThirdPartyService: 500,
	CodeRedisError:        500,
	CodeEmailServiceError: 500,
	CodeSMSServiceError:   500,
	CodePaymentError:      500,

	CodeFileTooLarge:            400,
	CodeFileTypeNotSupported:    400,
	CodeFileUploadFailed:        500,
	CodeFileNotFound:            404,
	CodeFileDeleteFailed:        500,
	CodeFileUpdateFailed:        500,
	CodeFileDownloadFailed:      500,
	CodeFileAccessDenied:        403,
	CodeStorageLimitExceeded:    400,
	CodeBandwidthLimitExceeded:  429,
	CodeUploadLimitExceeded:     400,
	CodeFileNotOwned:            403,
	CodeFileFormatNotSupport:    415,
	CodeStorageProviderNotFound: 404,

	CodeFolderNotFound:      404,
	CodeFolderCreateFailed:  500,
	CodeFolderDeleteFailed:  500,
	CodeFolderUpdateFailed:  500,
	CodeFolderNotEmpty:      400,
	CodeFolderNameDuplicate: 409,
	CodeInvalidFolderName:   400,

	CodeIPAccessDenied:       451,
	CodeIPNotInWhitelist:     451,
	CodeIPInBlacklist:        451,
	CodeRefererAccessDenied:  451,
	CodeDomainNotInWhitelist: 451,
	CodeDomainInBlacklist:    451,

	CodeSystemNotInstalled: 200,
	CodeInstallationFailed: 500,
	CodeConfigurationError: 500,

	CodeShareExpired:        400,
	CodeShareAccessExceeded: 403,
	CodeShareNotFound:       404,
	CodeSharePasswordWrong:  401,
}

var errorCodeToMessage = map[ErrorCode]string{
	CodeUnknown:            "服务器遇到了未知错误",
	CodeInternal:           "服务器内部错误",
	CodeInvalidParameter:   "无效的请求参数",
	CodeInvalidRequest:     "无效的请求",
	CodeUnauthorized:       "请先登录",
	CodeForbidden:          "没有操作权限",
	CodeNotFound:           "请求的资源不存在",
	CodeMethodNotAllowed:   "不支持的请求方法",
	CodeTimeout:            "请求超时",
	CodeConflict:           "资源冲突",
	CodeRateLimited:        "请求频率过高，请稍后再试",
	CodeValidationFailed:   "数据验证失败",
	CodeServiceUnavailable: "服务暂时不可用，请稍后再试",

	CodeUserNotFound:      "用户不存在",
	CodeWrongPassword:     "密码错误",
	CodeUserDisabled:      "账号已被禁用",
	CodeUserExists:        "用户已存在",
	CodeInvalidAuthToken:  "登录凭证无效，请重新登录",
	CodeExpiredAuthToken:  "登录已过期，请重新登录",
	CodeInvalidVerifyCode: "验证码无效或已过期",
	CodeEmailExists:       "邮箱已被注册",
	CodeEmailSendFailed:   "邮件发送失败，请稍后再试",
	CodeInvalidToken:      "无效的令牌",
	CodeTokenExpired:      "令牌已过期",
	CodeTokenUsed:         "令牌已被使用",

	CodeDBConnectionFailed: "数据库连接失败",
	CodeDBQueryFailed:      "数据查询失败",
	CodeDBCreateFailed:     "数据创建失败",
	CodeDBUpdateFailed:     "数据更新失败",
	CodeDBDeleteFailed:     "数据删除失败",
	CodeDBNoRecord:         "记录不存在",
	CodeDBDuplicate:        "记录已存在",
	CodeDBTransaction:      "数据库事务错误",
	CodeDBCommitFailed:     "提交数据库事务失败",

	CodeThirdPartyService: "第三方服务异常",
	CodeRedisError:        "缓存服务异常",
	CodeEmailServiceError: "邮件服务异常",
	CodeSMSServiceError:   "短信服务异常",
	CodePaymentError:      "支付服务异常",

	CodeFileTooLarge:            "文件文件过大",
	CodeFileTypeNotSupported:    "不支持的文件类型",
	CodeFileUploadFailed:        "文件上传失败",
	CodeFileNotFound:            "文件不存在",
	CodeFileDeleteFailed:        "文件删除失败",
	CodeFileUpdateFailed:        "文件更新失败",
	CodeFileDownloadFailed:      "文件下载失败",
	CodeFileAccessDenied:        "无权访问此文件",
	CodeStorageLimitExceeded:    "存储容量已用尽",
	CodeBandwidthLimitExceeded:  "带宽流量已用尽",
	CodeUploadLimitExceeded:     "上传次数已用尽",
	CodeFileNotOwned:            "文件不属于当前用户",
	CodeFileFormatNotSupport:    "文件格式不支持",
	CodeStorageProviderNotFound: "存储提供者未找到",

	CodeFolderNotFound:      "文件夹不存在",
	CodeFolderCreateFailed:  "文件夹创建失败",
	CodeFolderDeleteFailed:  "文件夹删除失败",
	CodeFolderUpdateFailed:  "文件夹更新失败",
	CodeFolderNotEmpty:      "文件夹不为空，无法删除",
	CodeFolderNameDuplicate: "文件夹名称已存在",
	CodeInvalidFolderName:   "无效的文件夹名称",

	CodeIPAccessDenied:       "IP访问已被禁止",
	CodeIPNotInWhitelist:     "您的IP不在访问白名单中",
	CodeIPInBlacklist:        "您的IP已被列入黑名单",
	CodeRefererAccessDenied:  "非法的请求来源",
	CodeDomainNotInWhitelist: "您的域名不在访问白名单中",
	CodeDomainInBlacklist:    "您的域名已被列入黑名单",

	CodeSystemNotInstalled: "系统未安装，请先完成安装配置",
	CodeInstallationFailed: "系统安装失败，请检查配置",
	CodeConfigurationError: "系统配置错误，请重新配置",

	CodeShareExpired:        "分享已过期",
	CodeShareAccessExceeded: "分享访问次数已超限",
	CodeShareNotFound:       "分享不存在",
	CodeSharePasswordWrong:  "分享密码错误",
}

type Error struct {
	Code      ErrorCode
	Message   string
	Detail    string
	Stack     string
	Time      time.Time
	RequestID string
	Metadata  map[string]interface{}
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func New(code ErrorCode, detail string) *Error {
	message, ok := errorCodeToMessage[code]
	if !ok {
		message = "未知错误"
	}
	var finalMessage string
	if message == detail {
		finalMessage = message
	} else {
		finalMessage = detail
	}
	err := &Error{
		Code:     code,
		Message:  finalMessage,
		Detail:   detail,
		Time:     time.Now(),
		Metadata: make(map[string]interface{}),
	}
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	var stackBuilder strings.Builder
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			fmt.Fprintf(&stackBuilder, "%s:%d - %s\n", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	err.Stack = stackBuilder.String()
	logger.Error("%s\nStack: %s", err.Error(), err.Stack)
	return err
}

func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

func (e *Error) WithMetadata(key string, value interface{}) *Error {
	e.Metadata[key] = value
	return e
}

func Is(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

func HTTPStatus(err error) int {
	if err == nil {
		return 200
	}
	if e, ok := err.(*Error); ok {
		if status, exists := errorCodeToHTTPStatus[e.Code]; exists {
			return status
		}
	}
	return 500
}

func Wrap(err error, code ErrorCode, detail string) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	e := New(code, detail)
	e.Detail = fmt.Sprintf("%s: %s", detail, err.Error())
	return e
}

func NewValidationError(field, detail string) *Error {
	err := New(CodeValidationFailed, detail)
	err.WithMetadata("field", field)
	return err
}

func GetSafeError(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return &Error{
			Code:      e.Code,
			Message:   e.Message,
			Time:      e.Time,
			RequestID: e.RequestID,
		}
	}
	return &Error{
		Code:    CodeInternal,
		Message: errorCodeToMessage[CodeInternal],
		Time:    time.Now(),
	}
}
