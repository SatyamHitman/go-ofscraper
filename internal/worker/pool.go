// =============================================================================
// FILE: internal/worker/pool.go
// PURPOSE: Generic worker pool. Provides a reusable concurrent worker pool
//          pattern using channels and goroutines. Used for download workers,
//          API pagination workers, and other concurrent batch operations.
// =============================================================================

package worker

import (
	"context"
	"sync"
)

// ---------------------------------------------------------------------------
// Job and Result types
// ---------------------------------------------------------------------------

// Job represents a unit of work to be processed by a worker.
type Job[T any] struct {
	ID    int
	Input T
}

// Result represents the outcome of processing a job.
type Result[T any, R any] struct {
	JobID  int
	Input  T
	Output R
	Err    error
}

// ---------------------------------------------------------------------------
// Pool
// ---------------------------------------------------------------------------

// Pool is a generic worker pool that processes jobs concurrently.
type Pool[T any, R any] struct {
	workers int
	handler func(context.Context, T) (R, error)
}

// NewPool creates a worker pool with the given concurrency and handler function.
//
// Parameters:
//   - workers: Number of concurrent workers.
//   - handler: Function that processes a single job input and returns a result.
//
// Returns:
//   - A configured Pool.
func NewPool[T any, R any](workers int, handler func(context.Context, T) (R, error)) *Pool[T, R] {
	if workers <= 0 {
		workers = 1
	}
	return &Pool[T, R]{
		workers: workers,
		handler: handler,
	}
}

// Run processes all inputs concurrently and returns results.
// Results are returned in no guaranteed order.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - inputs: Items to process.
//
// Returns:
//   - Slice of Results.
func (p *Pool[T, R]) Run(ctx context.Context, inputs []T) []Result[T, R] {
	jobs := make(chan Job[T], len(inputs))
	results := make(chan Result[T, R], len(inputs))

	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if ctx.Err() != nil {
					results <- Result[T, R]{
						JobID: job.ID,
						Input: job.Input,
						Err:   ctx.Err(),
					}
					continue
				}

				output, err := p.handler(ctx, job.Input)
				results <- Result[T, R]{
					JobID:  job.ID,
					Input:  job.Input,
					Output: output,
					Err:    err,
				}
			}
		}()
	}

	// Enqueue jobs.
	for i, input := range inputs {
		jobs <- Job[T]{ID: i, Input: input}
	}
	close(jobs)

	// Wait for completion and collect results.
	go func() {
		wg.Wait()
		close(results)
	}()

	var collected []Result[T, R]
	for r := range results {
		collected = append(collected, r)
	}
	return collected
}

// ---------------------------------------------------------------------------
// Simple pool (non-generic, for when T = R = interface{})
// ---------------------------------------------------------------------------

// SimplePool runs a batch of functions concurrently.
type SimplePool struct {
	workers int
}

// NewSimplePool creates a simple worker pool.
func NewSimplePool(workers int) *SimplePool {
	if workers <= 0 {
		workers = 1
	}
	return &SimplePool{workers: workers}
}

// Run executes all functions concurrently with bounded parallelism.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - tasks: Functions to execute.
//
// Returns:
//   - Slice of errors (nil for successful tasks), in order.
func (sp *SimplePool) Run(ctx context.Context, tasks []func(context.Context) error) []error {
	errs := make([]error, len(tasks))
	sem := make(chan struct{}, sp.workers)
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, fn func(context.Context) error) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				errs[idx] = ctx.Err()
				return
			}
			errs[idx] = fn(ctx)
		}(i, task)
	}

	wg.Wait()
	return errs
}
