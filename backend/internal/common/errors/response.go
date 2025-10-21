package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResultVO defines the standard API response structure.
type ResultVO struct {
	Code int         `json:"code"`           // 状态码
	Msg  string      `json:"msg"`            // 提示消息
	Data interface{} `json:"data,omitempty"` // 返回数据
}

const (
	// SuccessCode represents the code for a successful operation.
	SuccessCode = 200
	// FailCode represents the code for a failed operation.
	FailCode = 400
)

// Success sends a successful response with a message and data.
func Success(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(ResultVO{
		Code: SuccessCode,
		Msg:  msg,
		Data: data,
	})
	if err != nil {
		log.Printf("failed to write success response: %v", err)
	}
}

// SuccessWithData sends a successful response with default message and data.
func SuccessWithData(w http.ResponseWriter, data interface{}) {
	Success(w, "ok", data)
}

// SuccessWithoutData sends a successful response with a message and no data.
func SuccessWithoutData(w http.ResponseWriter, msg string) {
	Success(w, msg, nil)
}

// Ok sends a successful response with a default message and no data.
func Ok(w http.ResponseWriter) {
	Success(w, "ok", nil)
}

// Fail sends a failure response with a specific HTTP status code, business code, and message.
func Fail(w http.ResponseWriter, httpStatusCode, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(ResultVO{
		Code: code,
		Msg:  msg,
	})
	if err != nil {
		log.Printf("failed to write fail response: %v", err)
	}
}

// FailWithMessage sends a failure response with a default HTTP 400 status code.
func FailWithMessage(w http.ResponseWriter, msg string) {
	Fail(w, http.StatusBadRequest, FailCode, msg)
}

// FailWithCode sends a failure response with a specific business code and message.
func FailWithCode(w http.ResponseWriter, code int, msg string) {
	Fail(w, http.StatusBadRequest, code, msg)
}
