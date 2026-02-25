// =============================================================================
// FILE: internal/tui/live/tasks.go
// PURPOSE: Task display. Shows running tasks with spinners and status
//          indicators. Ports Python utils/live/tasks.py.
// =============================================================================

package live

import (
	"fmt"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Task display
// ---------------------------------------------------------------------------

// TaskStatus represents the state of a task.
type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskRunning
	TaskDone
	TaskFailed
)

// Task represents a displayable task.
type Task struct {
	Name      string
	Status    TaskStatus
	StartTime time.Time
	EndTime   time.Time
	Message   string
}

// TaskDisplay manages a list of displayable tasks.
type TaskDisplay struct {
	mu    sync.Mutex
	tasks []*Task
}

// NewTaskDisplay creates a new task display.
func NewTaskDisplay() *TaskDisplay {
	return &TaskDisplay{}
}

// AddTask adds a new task to the display.
func (td *TaskDisplay) AddTask(name string) *Task {
	td.mu.Lock()
	defer td.mu.Unlock()
	t := &Task{
		Name:   name,
		Status: TaskPending,
	}
	td.tasks = append(td.tasks, t)
	return t
}

// StartTask marks a task as running.
func (td *TaskDisplay) StartTask(t *Task) {
	td.mu.Lock()
	defer td.mu.Unlock()
	t.Status = TaskRunning
	t.StartTime = time.Now()
}

// CompleteTask marks a task as done.
func (td *TaskDisplay) CompleteTask(t *Task, msg string) {
	td.mu.Lock()
	defer td.mu.Unlock()
	t.Status = TaskDone
	t.EndTime = time.Now()
	t.Message = msg
}

// FailTask marks a task as failed.
func (td *TaskDisplay) FailTask(t *Task, msg string) {
	td.mu.Lock()
	defer td.mu.Unlock()
	t.Status = TaskFailed
	t.EndTime = time.Now()
	t.Message = msg
}

// Render returns the task list as displayable lines.
func (td *TaskDisplay) Render() []string {
	td.mu.Lock()
	defer td.mu.Unlock()

	spinChars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	spinIdx := int(time.Now().UnixMilli()/100) % len(spinChars)

	var lines []string
	for _, t := range td.tasks {
		var prefix string
		switch t.Status {
		case TaskPending:
			prefix = "○"
		case TaskRunning:
			prefix = string(spinChars[spinIdx])
		case TaskDone:
			prefix = "✓"
		case TaskFailed:
			prefix = "✗"
		}

		line := fmt.Sprintf(" %s %s", prefix, t.Name)
		if t.Message != "" {
			line += fmt.Sprintf(" — %s", t.Message)
		}
		if t.Status == TaskRunning && !t.StartTime.IsZero() {
			elapsed := time.Since(t.StartTime).Round(time.Second)
			line += fmt.Sprintf(" (%s)", elapsed)
		}
		lines = append(lines, line)
	}
	return lines
}
