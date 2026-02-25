// =============================================================================
// FILE: internal/tui/live_classes/task.go
// PURPOSE: Task status data class. Stores task name, current status, and
//          elapsed time for display in the live TUI.
//          Ports Python utils/live/classes/task.py.
// =============================================================================

package liveclasses

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// TaskStatus
// ---------------------------------------------------------------------------

// TaskState represents the current state of a task.
type TaskState int

const (
	TaskPending TaskState = iota
	TaskRunning
	TaskCompleted
	TaskFailed
	TaskSkipped
)

// taskStateNames maps state to display string.
var taskStateNames = []string{"Pending", "Running", "Completed", "Failed", "Skipped"}

// String returns the human-readable name for the task state.
func (s TaskState) String() string {
	if int(s) < len(taskStateNames) {
		return taskStateNames[s]
	}
	return "Unknown"
}

// ---------------------------------------------------------------------------
// TaskInfo
// ---------------------------------------------------------------------------

// TaskInfo holds status data for a single task in the live display.
type TaskInfo struct {
	// Identification.
	ID   string
	Name string

	// Status.
	State   TaskState
	Message string

	// Timing.
	StartTime time.Time
	EndTime   time.Time

	// Progress (optional, 0-100).
	Progress float64

	// Sub-tasks (optional).
	SubTasks []*TaskInfo
}

// NewTaskInfo creates a new TaskInfo in the pending state.
func NewTaskInfo(id, name string) *TaskInfo {
	return &TaskInfo{
		ID:    id,
		Name:  name,
		State: TaskPending,
	}
}

// Start marks the task as running.
func (t *TaskInfo) Start() {
	t.State = TaskRunning
	t.StartTime = time.Now()
}

// Complete marks the task as completed with an optional message.
func (t *TaskInfo) Complete(msg string) {
	t.State = TaskCompleted
	t.Message = msg
	t.EndTime = time.Now()
}

// Fail marks the task as failed with an error message.
func (t *TaskInfo) Fail(msg string) {
	t.State = TaskFailed
	t.Message = msg
	t.EndTime = time.Now()
}

// Skip marks the task as skipped.
func (t *TaskInfo) Skip(msg string) {
	t.State = TaskSkipped
	t.Message = msg
	t.EndTime = time.Now()
}

// SetProgress sets the task progress percentage (0-100).
func (t *TaskInfo) SetProgress(pct float64) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	t.Progress = pct
}

// Elapsed returns the duration since the task started.
// Returns zero if the task has not started.
func (t *TaskInfo) Elapsed() time.Duration {
	if t.StartTime.IsZero() {
		return 0
	}
	end := t.EndTime
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(t.StartTime).Round(time.Millisecond)
}

// IsActive returns true if the task is pending or running.
func (t *TaskInfo) IsActive() bool {
	return t.State == TaskPending || t.State == TaskRunning
}

// IsDone returns true if the task has completed, failed, or been skipped.
func (t *TaskInfo) IsDone() bool {
	return t.State == TaskCompleted || t.State == TaskFailed || t.State == TaskSkipped
}

// AddSubTask adds a sub-task.
func (t *TaskInfo) AddSubTask(sub *TaskInfo) {
	t.SubTasks = append(t.SubTasks, sub)
}

// StatusLine returns a single-line summary of the task for display.
func (t *TaskInfo) StatusLine() string {
	prefix := " "
	switch t.State {
	case TaskPending:
		prefix = "o"
	case TaskRunning:
		prefix = ">"
	case TaskCompleted:
		prefix = "+"
	case TaskFailed:
		prefix = "x"
	case TaskSkipped:
		prefix = "-"
	}

	line := fmt.Sprintf("[%s] %s: %s", prefix, t.Name, t.State)
	if t.Message != "" {
		line += fmt.Sprintf(" - %s", t.Message)
	}
	if t.State == TaskRunning {
		line += fmt.Sprintf(" (%s)", t.Elapsed().Round(time.Second))
	}
	return line
}
