package common

import (
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogWrapperConsumer wraps a consumer with logging functionality
type LogWrapperConsumer struct {
	target       func(*RunningStatus) // Target consumer function
	finished     *bool                // Last finished status
	lastTask     string               // Last task name
	lastProgress *Float3              // Last progress value
}

// NewLogWrapperConsumer creates a new LogWrapperConsumer
func NewLogWrapperConsumer(target func(*RunningStatus)) *LogWrapperConsumer {
	return &LogWrapperConsumer{
		target: target,
	}
}

// Accept processes the running status with logging
func (l *LogWrapperConsumer) Accept(status *RunningStatus) {
	// Log the status as JSON
	jsonData, _ := json.Marshal(status)
	logx.Infof("** status: %s", string(jsonData))

	// Update internal state
	l.finished = status.Finished
	if status.Task != "" {
		l.lastTask = status.Task
	}
	if progress := status.Progress; progress != nil {
		l.lastProgress = progress
	}

	// Call target consumer
	if l.target != nil {
		l.target(status)
	}
}

// GetFinished returns the last finished status
func (l *LogWrapperConsumer) GetFinished() *bool {
	return l.finished
}

// GetLastTask returns the last task name
func (l *LogWrapperConsumer) GetLastTask() string {
	return l.lastTask
}

// GetLastProgress returns the last progress value
func (l *LogWrapperConsumer) GetLastProgress() *Float3 {
	return l.lastProgress
}
