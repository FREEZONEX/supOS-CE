package common

import (
	"encoding/json"
	"log"
)

// LogWrapperConsumer wraps a consumer with logging functionality
// Corresponds to Java's LogWrapperConsumer.java
type LogWrapperConsumer struct {
	target       func(*RunningStatus) // Target consumer function
	finished     *bool                // Last finished status
	lastTask     string               // Last task name
	lastProgress *float64             // Last progress value
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
	log.Printf("** status: %s", string(jsonData))

	// Update internal state
	l.finished = status.Finished
	if status.Task != "" {
		l.lastTask = status.Task
	}
	if status.Progress != nil {
		l.lastProgress = status.Progress
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
func (l *LogWrapperConsumer) GetLastProgress() *float64 {
	return l.lastProgress
}
