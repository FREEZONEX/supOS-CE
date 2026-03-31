package errors

import (
	"backend/internal/common/I18nUtils"
	"context"
	"fmt"
	"strings"
)

// BuzError represents a business logic error
type BuzError struct {
	Code int
	Msg  string
}

// Error implements the error interface
func (e *BuzError) Error() string {
	return fmt.Sprintf(`{"code":%d, "msg":"%s"}`, e.Code, jsonEscape(e.Msg))
}

// jsonEscape escapes special characters in a string for safe JSON encoding
func jsonEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)  // 反斜杠
	s = strings.ReplaceAll(s, `"`, `\"`)  // 双引号
	s = strings.ReplaceAll(s, "\n", `\n`) // 换行符
	s = strings.ReplaceAll(s, "\r", `\r`) // 回车符
	s = strings.ReplaceAll(s, "\t", `\t`) // 制表符
	s = strings.ReplaceAll(s, "\b", `\b`) // 退格符
	s = strings.ReplaceAll(s, "\f", `\f`) // 换页符
	return s
}

// NewBuzError creates a new business error with code and message
func NewBuzError(ctx context.Context, code int, msg string, params ...any) *BuzError {
	return &BuzError{
		Code: code,
		Msg:  I18nUtils.GetMessageWithCtx(ctx, msg, params...),
	}
}

// NewBuzErrorWithMsg creates a new business error with message only
func NewBuzErrorWithMsg(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 500, msg, params...)
}

// AppError represents an application error
type AppError struct {
	*BuzError
}

// NewAppError creates a new application error with code and message
func NewAppError(ctx context.Context, code int, msg string) *AppError {
	return &AppError{
		BuzError: NewBuzError(ctx, code, msg),
	}
}

// NewAppErrorWithMsg creates a new application error with message only
func NewAppErrorWithMsg(ctx context.Context, msg string) *AppError {
	return &AppError{
		BuzError: NewBuzErrorWithMsg(ctx, msg),
	}
}

// NodeRedError represents a Node-RED adapter error
type NodeRedError struct {
	*BuzError
}

// NewNodeRedError creates a new Node-RED error with code and message
func NewNodeRedError(ctx context.Context, code int, msg string, params ...any) *NodeRedError {
	return &NodeRedError{
		BuzError: NewBuzError(ctx, code, msg, params...),
	}
}

// NewNodeRedErrorWithMsg creates a new Node-RED error with message only
func NewNodeRedErrorWithMsg(ctx context.Context, msg string, params ...any) *NodeRedError {
	return &NodeRedError{
		BuzError: NewBuzErrorWithMsg(ctx, msg, params...),
	}
}

// Common error constructors

// BadRequest creates a 400 bad request error
func BadRequest(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 400, msg, params...)
}

// Unauthorized creates a 401 unauthorized error
func Unauthorized(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 401, msg, params...)
}

// Forbidden creates a 403 forbidden error
func Forbidden(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 403, msg, params...)
}

// NotFound creates a 404 not found error
func NotFound(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 404, msg, params...)
}

// InternalError creates a 500 internal server error
func InternalError(ctx context.Context, msg string, params ...any) *BuzError {
	return NewBuzError(ctx, 500, msg, params...)
}

// IsBuzError checks if error is a BuzError
func IsBuzError(err error) (*BuzError, bool) {
	if buzErr, ok := err.(*BuzError); ok {
		return buzErr, true
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr.BuzError, true
	}
	if nodeRedErr, ok := err.(*NodeRedError); ok {
		return nodeRedErr.BuzError, true
	}
	return nil, false
}

const (
	// ... other error codes
	UserAlreadyExists = 20001
)
