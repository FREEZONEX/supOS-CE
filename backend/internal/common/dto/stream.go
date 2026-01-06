package dto

import (
	"backend/internal/common/constants"
	"fmt"
	"strings"
)

// Stream status constants
const (
	StreamStatusUnknown = 0 // 未知
	StreamStatusRunning = 1 // 在运行
	StreamStatusPaused  = 2 // 已暂停
)

// StreamInfo represents stream computation information
type StreamInfo struct {
	Name        string `json:"name"`        // 流计算名称
	SQL         string `json:"sql"`         // SQL query
	TargetTable string `json:"targetTable"` // 目标表
	Status      int    `json:"status"`      // 0--未知, 1--在运行, 2--已暂停
}

// NewStreamInfo creates a new StreamInfo
func NewStreamInfo(name, targetTable string, status int) *StreamInfo {
	return &StreamInfo{
		Name:        name,
		TargetTable: targetTable,
		Status:      status,
	}
}

// StreamOptions represents stream options
type StreamOptions struct {
	// Base stream options fields
}

// StreamWindowOptions represents stream window options
type StreamWindowOptions struct {
	// Base window options fields
}

// StreamWindowOptionsSession represents session window options
type StreamWindowOptionsSession struct {
	TolValue string `json:"tolValue" validate:"required"` // 会话超时时间
}

// String returns the string representation
func (s *StreamWindowOptionsSession) String() string {
	var builder strings.Builder
	builder.WriteString("SESSION(")
	builder.WriteString(constants.SysFieldCreateTime)
	builder.WriteString(",")
	builder.WriteString(s.TolValue)
	builder.WriteString(")")
	return builder.String()
}

// StreamWindowOptionsInterval represents interval window options
type StreamWindowOptionsInterval struct {
	IntervalValue  string  `json:"intervalValue" validate:"required"` // 时间间隔
	IntervalOffset *string `json:"intervalOffset,omitzero"`           // 时间偏移量
}

// String returns the string representation
func (s *StreamWindowOptionsInterval) String() string {
	var builder strings.Builder
	builder.WriteString("INTERVAL(")
	builder.WriteString(s.IntervalValue)
	if s.IntervalOffset != nil && *s.IntervalOffset != "" {
		builder.WriteString(",")
		builder.WriteString(*s.IntervalOffset)
	}
	builder.WriteString(")")
	return builder.String()
}

// StreamWindowOptionsEventWindow represents event window options
type StreamWindowOptionsEventWindow struct {
	StartWith string `json:"startWith" validate:"required"` // 开始条件
	EndWith   string `json:"endWith" validate:"required"`   // 结束条件
}

// String returns the string representation
func (s *StreamWindowOptionsEventWindow) String() string {
	var builder strings.Builder
	builder.WriteString("EVENT_WINDOW START WITH ")
	builder.WriteString(s.StartWith)
	builder.WriteString(" END WITH ")
	builder.WriteString(s.EndWith)
	return builder.String()
}

// StreamWindowOptionsStateWindow represents state window options
type StreamWindowOptionsStateWindow struct {
	Field string `json:"field" validate:"required"` // 状态字段
}

// String returns the string representation
func (s *StreamWindowOptionsStateWindow) String() string {
	return fmt.Sprintf("STATE_WINDOW(%s)", s.Field)
}

// StreamWindowOptionsCountWindow represents count window options
type StreamWindowOptionsCountWindow struct {
	CountValue   int  `json:"countValue" validate:"required,min=2"` // 计数值
	SlidingValue *int `json:"slidingValue,omitzero"`                // 滑动值
}

// String returns the string representation
func (s *StreamWindowOptionsCountWindow) String() string {
	var builder strings.Builder
	builder.WriteString("COUNT_WINDOW(")
	builder.WriteString(fmt.Sprintf("%d", s.CountValue))
	if s.SlidingValue != nil {
		builder.WriteString(",")
		builder.WriteString(fmt.Sprintf("%d", *s.SlidingValue))
	}
	builder.WriteString(")")
	return builder.String()
}
