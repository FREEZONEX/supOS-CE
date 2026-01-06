package dto

// DashboardExportParam Dashboard 导出参数
type DashboardExportParam struct {
	Ids        []string `json:"ids"`
	ExportType string   `json:"exportType"` // "ALL" or selected ids
}

// RunningStatus 异步任务运行状态
type RunningStatus struct {
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	Data         string  `json:"data,omitempty"` // 通常用于返回错误文件的路径
	Task         string  `json:"task,omitempty"`
	Progress     float64 `json:"progress,omitempty"`
	TotalCount   int     `json:"totalCount,omitempty"`
	SuccessCount int     `json:"successCount,omitempty"`
	ErrorCount   int     `json:"errorCount,omitempty"`
}

// NewRunningStatus 创建 RunningStatus 实例
func NewRunningStatus(code int, msg string, data ...string) *RunningStatus {
	status := &RunningStatus{
		Code: code,
		Msg:  msg,
	}
	if len(data) > 0 {
		status.Data = data[0]
	}
	return status
}

func (rs *RunningStatus) SetTask(task string) *RunningStatus {
	rs.Task = task
	return rs
}

func (rs *RunningStatus) SetProgress(progress float64) *RunningStatus {
	rs.Progress = progress
	return rs
}
