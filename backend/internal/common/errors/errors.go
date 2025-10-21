package errors

import "fmt"

// BuzError represents a business logic error
type BuzError struct {
	Code   int
	Msg    string
	Params []any
}

// Error implements the error interface
func (e *BuzError) Error() string {
	if len(e.Params) > 0 {
		return fmt.Sprintf(e.Msg, e.Params...)
	}
	return e.Msg
}

// NewBuzError creates a new business error with code and message
func NewBuzError(code int, msg string, params ...any) *BuzError {
	return &BuzError{
		Code:   code,
		Msg:    msg,
		Params: params,
	}
}

// NewBuzErrorWithMsg creates a new business error with message only
func NewBuzErrorWithMsg(msg string, params ...any) *BuzError {
	return &BuzError{
		Code:   500,
		Msg:    msg,
		Params: params,
	}
}

// AppError represents an application error
type AppError struct {
	*BuzError
}

// NewAppError creates a new application error with code and message
func NewAppError(code int, msg string) *AppError {
	return &AppError{
		BuzError: NewBuzError(code, msg),
	}
}

// NewAppErrorWithMsg creates a new application error with message only
func NewAppErrorWithMsg(msg string) *AppError {
	return &AppError{
		BuzError: NewBuzErrorWithMsg(msg),
	}
}

// NodeRedError represents a Node-RED adapter error
type NodeRedError struct {
	*BuzError
}

// NewNodeRedError creates a new Node-RED error with code and message
func NewNodeRedError(code int, msg string, params ...any) *NodeRedError {
	return &NodeRedError{
		BuzError: NewBuzError(code, msg, params...),
	}
}

// NewNodeRedErrorWithMsg creates a new Node-RED error with message only
func NewNodeRedErrorWithMsg(msg string, params ...any) *NodeRedError {
	return &NodeRedError{
		BuzError: NewBuzErrorWithMsg(msg, params...),
	}
}

// Common error constructors

// BadRequest creates a 400 bad request error
func BadRequest(msg string, params ...any) *BuzError {
	return NewBuzError(400, msg, params...)
}

// Unauthorized creates a 401 unauthorized error
func Unauthorized(msg string, params ...any) *BuzError {
	return NewBuzError(401, msg, params...)
}

// Forbidden creates a 403 forbidden error
func Forbidden(msg string, params ...any) *BuzError {
	return NewBuzError(403, msg, params...)
}

// NotFound creates a 404 not found error
func NotFound(msg string, params ...any) *BuzError {
	return NewBuzError(404, msg, params...)
}

// InternalError creates a 500 internal server error
func InternalError(msg string, params ...any) *BuzError {
	return NewBuzError(500, msg, params...)
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
