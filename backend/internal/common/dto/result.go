package dto

// BaseResult represents basic API response structure
type BaseResult struct {
	Code int    `json:"code"` // 错误码，0--正常，其他失败
	Msg  string `json:"msg"`  // 错误信息
}

// ResultDTO represents generic API response with data
type ResultDTO[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitzero"`
	Data T      `json:"data,omitzero"`
}

// SuccessWithData creates a successful result with data
func SuccessWithData[T any](data T) ResultDTO[T] {
	return ResultDTO[T]{
		Code: 200,
		Data: data,
	}
}

// Success creates a successful result with message
func Success(msg string) ResultDTO[any] {
	return ResultDTO[any]{
		Code: 200,
		Msg:  msg,
	}
}

// Fail creates a failed result with message
func Fail(msg string) ResultDTO[any] {
	return ResultDTO[any]{
		Code: 500,
		Msg:  msg,
	}
}

// IsSuccess checks if the result is successful
func (r *ResultDTO[T]) IsSuccess() bool {
	return r.Code == 200
}

// JsonResult represents JSON API response
type JsonResult[T any] struct {
	BaseResult
	Data T `json:"data,omitzero"`
}

// NewJsonResult creates a new JSON result
func NewJsonResult[T any](code int, msg string, data T) *JsonResult[T] {
	return &JsonResult[T]{
		BaseResult: BaseResult{
			Code: code,
			Msg:  msg,
		},
		Data: data,
	}
}
