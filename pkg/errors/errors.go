package errors

import "fmt"

// APIError API错误结构
type APIError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现错误包装接口
func (e *APIError) Unwrap() error {
	return e.Err
}

// New 创建新的API错误
func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装现有错误
func Wrap(err error, code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common errors 常见错误定义
var (
	ErrInternalServer = &APIError{Code: 500, Message: "服务器内部错误"}
	ErrBadRequest     = &APIError{Code: 400, Message: "请求参数错误"}
	ErrUnauthorized   = &APIError{Code: 401, Message: "未授权访问"}
	ErrForbidden      = &APIError{Code: 403, Message: "禁止访问"}
	ErrNotFound       = &APIError{Code: 404, Message: "资源不存在"}
	ErrValidateFailed = &APIError{Code: 422, Message: "参数校验失败"}
)
