// =============================================================================
// FILE: internal/commands/scraper/jobqueue.go
// PURPOSE: Job queue for managing scrape jobs with priority and concurrent
//          processing. Supports enqueuing jobs for users and dispatching them
//          to workers. Ports Python runner/manager/batch_manager.py queue logic.
// =============================================================================

package scraper

import (
	"container/heap"
	"context"
	"log/slog"
	"sync"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Job represents a single scrape job for a user.
// ---------------------------------------------------------------------------

// Job holds the data needed to process a single user's scrape.
type Job struct {
	User     *model.User
	Areas    []string
	Actions  []string
	Priority int // Lower number = higher priority.
	index    int // Heap index, managed by the priority queue.
}

// ---------------------------------------------------------------------------
// jobHeap implements heap.Interface for priority-based job ordering.
// ---------------------------------------------------------------------------

type jobHeap []*Job

func (h jobHeap) Len() int           { return len(h) }
func (h jobHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }
func (h jobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *jobHeap) Push(x interface{}) {
	n := len(*h)
	job := x.(*Job)
	job.index = n
	*h = append(*h, job)
}

func (h *jobHeap) Pop() interface{} {
	old := *h
	n := len(old)
	job := old[n-1]
	old[n-1] = nil
	job.index = -1
	*h = old[:n-1]
	return job
}

// ---------------------------------------------------------------------------
// JobQueue manages queued scrape jobs with priority ordering.
// ---------------------------------------------------------------------------

// JobQueue provides a thread-safe priority queue for scrape jobs.
type JobQueue struct {
	mu     sync.Mutex
	jobs   jobHeap
	logger *slog.Logger
}

// NewJobQueue creates a new empty job queue.
//
// Parameters:
//   - logger: Structured logger for queue operations.
//
// Returns:
//   - An initialized JobQueue.
func NewJobQueue(logger *slog.Logger) *JobQueue {
	if logger == nil {
		logger = slog.Default()
	}
	jq := &JobQueue{
		logger: logger,
	}
	heap.Init(&jq.jobs)
	return jq
}

// Enqueue adds a job to the queue.
//
// Parameters:
//   - job: The scrape job to enqueue.
func (jq *JobQueue) Enqueue(job *Job) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	heap.Push(&jq.jobs, job)
	jq.logger.Debug("job enqueued",
		"user", job.User.Name,
		"priority", job.Priority,
		"queue_size", jq.jobs.Len(),
	)
}

// Dequeue removes and returns the highest-priority job from the queue.
//
// Returns:
//   - The next Job, or nil if the queue is empty.
func (jq *JobQueue) Dequeue() *Job {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	if jq.jobs.Len() == 0 {
		return nil
	}
	return heap.Pop(&jq.jobs).(*Job)
}

// Len returns the number of pending jobs in the queue.
//
// Returns:
//   - The queue length.
func (jq *JobQueue) Len() int {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	return jq.jobs.Len()
}

// ProcessAll dequeues and processes all jobs using the provided handler function.
// Runs up to maxWorkers jobs concurrently.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - maxWorkers: Maximum number of concurrent workers.
//   - handler: Function to process each job.
//
// Returns:
//   - The number of jobs processed, and the first error encountered.
func (jq *JobQueue) ProcessAll(ctx context.Context, maxWorkers int, handler func(ctx context.Context, job *Job) error) (int, error) {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	processed := 0

	for {
		job := jq.Dequeue()
		if job == nil {
			break
		}

		if ctx.Err() != nil {
			return processed, ctx.Err()
		}

		sem <- struct{}{}
		wg.Add(1)
		processed++

		go func(j *Job) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := handler(ctx, j); err != nil {
				jq.logger.Error("job failed",
					"user", j.User.Name,
					"error", err,
				)
				errOnce.Do(func() { firstErr = err })
			}
		}(job)
	}

	wg.Wait()
	return processed, firstErr
}
