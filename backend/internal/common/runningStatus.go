package common

import (
	"backend/internal/common/I18nUtils"
	"strconv"
)

// RunningStatus represents running status information
type RunningStatus struct {
	Module     string  `json:"module,omitempty"`     // 模块 uns sourceFlow eventFlow dashboard
	Code       int     `json:"code"`                 // 状态码 200表示成功
	Msg        string  `json:"msg,omitempty"`        // 消息
	ErrTipFile string  `json:"errTipFile,omitempty"` // 错误提示文件
	N          *int    `json:"n,omitempty"`          // 总数
	I          *int    `json:"i,omitempty"`          // 当前索引
	Task       string  `json:"task,omitempty"`       // 任务名称
	SpendMills *int64  `json:"spendMills,omitempty"` // 耗时毫秒
	Finished   *bool   `json:"finished,omitempty"`   // 是否完成
	Progress   *Float3 `json:"progress,omitempty"`   // 进度：0-100

	StartTime    *int64 `json:"startTime,omitempty"`    // 开始时间
	EndTime      *int64 `json:"endTime,omitempty"`      // 结束时间
	TotalCount   int    `json:"totalCount,omitempty"`   // 总数量
	ErrorCount   int    `json:"errorCount,omitempty"`   // 错误数量
	SuccessCount int    `json:"successCount,omitempty"` // 成功数量
}

// 自定义浮点类型，保留3位小数
type Float3 float64

func (f Float3) MarshalJSON() ([]byte, error) {
	// 格式化为保留3位小数的字符串
	formatted := strconv.FormatFloat(float64(f), 'f', 3, 64)
	return []byte(formatted), nil
}

// NewRunningStatus creates a new RunningStatus with default values
func NewRunningStatus() *RunningStatus {
	return &RunningStatus{}
}

// NewRunningStatusWithCode creates a new RunningStatus with code and message
func NewRunningStatusWithCode(code int, msg string) *RunningStatus {
	finished := true
	return &RunningStatus{
		Code:     code,
		Msg:      I18nUtils.GetMessage(msg),
		Finished: &finished,
	}
}

// NewRunningStatusWithError creates a new RunningStatus with error information
func NewRunningStatusWithError(code int, msg string, errTipFile string) *RunningStatus {
	finished := true
	return &RunningStatus{
		Code:       code,
		Msg:        I18nUtils.GetMessage(msg),
		ErrTipFile: errTipFile,
		Finished:   &finished,
	}
}

// NewRunningStatusWithProgress creates a new RunningStatus with progress information
func NewRunningStatusWithProgress(n, i int, task, msg string) *RunningStatus {
	return &RunningStatus{
		N:    &n,
		I:    &i,
		Task: task,
		Msg:  I18nUtils.GetMessage(msg),
	}
}

// SetSpendMills sets the spend time and calculates progress
func (r *RunningStatus) SetSpendMills(spend int64) *RunningStatus {
	r.SpendMills = &spend
	if r.N != nil && *r.N > 0 && r.I != nil {
		progress := Float3(1000*float64(*r.I)/float64(*r.N)) / 10.0
		r.Progress = &progress
	}
	return r
}

// SetProgress sets the progress value
func (r *RunningStatus) SetProgress(progress float64) *RunningStatus {
	prog := Float3(progress)
	r.Progress = &prog
	return r
}

// SetCode sets the code
func (r *RunningStatus) SetCode(code int) *RunningStatus {
	r.Code = code
	return r
}

// SetTask sets the task name
func (r *RunningStatus) SetTask(task string) *RunningStatus {
	r.Task = task
	return r
}

// SetStartTime sets the start time
func (r *RunningStatus) SetStartTime(startTime int64) *RunningStatus {
	r.StartTime = &startTime
	return r
}

// SetEndTime sets the end time
func (r *RunningStatus) SetEndTime(endTime int64) *RunningStatus {
	r.EndTime = &endTime
	return r
}

// SetFinished sets the finished status
func (r *RunningStatus) SetFinished(finished bool) *RunningStatus {
	r.Finished = &finished
	return r
}
