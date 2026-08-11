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
}

func (e *HTTPError) Error() string { return e.Msg }

func NewHTTPError(status int, msg string) error {
	return &HTTPError{Status: status, Code: status, Msg: msg}
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
