package httpx

import (
	"errors"
	"net/http"

	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type HTTPError struct {
	Status int
	Code   int
	Msg    string
	// cause 保留原始错误，供 Unwrap 透出；nil 表示无原始错误。
	cause error
}

func (e *HTTPError) Error() string { return e.Msg }

// Unwrap 暴露原始错误，使 errors.Is/errors.As 能穿透 HTTPError 找到业务错误
// （validateContainerAppSave 等 logic 层直接返回的包装错误可被测试与调用方判别）。
func (e *HTTPError) Unwrap() error { return e.cause }

func NewHTTPError(status int, msg string) error {
	return &HTTPError{Status: status, Code: status, Msg: msg}
}

// WrapHTTPError 生成携带原始错误的 HTTPError，errors.Is 可穿透到 cause。
func WrapHTTPError(status int, msg string, cause error) error {
	return &HTTPError{Status: status, Code: status, Msg: msg, cause: cause}
}

func Envelope(data any) *types.Envelope {
	return &types.Envelope{Code: 200, Msg: "success", Data: data}
}

func Empty() *types.Envelope {
	return Envelope(map[string]any{})
}

func WriteError(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		httpx.WriteJson(w, httpErr.Status, &types.Envelope{Code: httpErr.Code, Msg: httpErr.Msg})
		return
	}
	httpx.WriteJson(w, http.StatusInternalServerError, &types.Envelope{Code: http.StatusInternalServerError, Msg: err.Error()})
}
