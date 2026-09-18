package worker

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Task struct {
	VideoID uuid.UUID
}

type WorkerPool struct {
	tasksChan chan Task
	wg        sync.WaitGroup
	enrichFn  func(ctx context.Context, videoID uuid.UUID) error
}

// a new pool
func NewWorkerPool(bufferSize int, enrichFn func(ctx context.Context, videoID uuid.UUID) error) *WorkerPool {
	return &WorkerPool{
		tasksChan: make(chan Task, bufferSize),
		enrichFn:  enrichFn,
	}
}

// spawn workers
// main() asks for 3
// the goroutines wait for items to hit the tasksChan
func (wp *WorkerPool) Start(ctx context.Context, numWorkers int) {
	for i := range numWorkers {
		wp.wg.Add(1)
		go wp.runWorker(ctx, i)
	}
}

// api calls enqueue
// enqueue attempts to push a task to the queue; returns false if full
// task -> tasksChan
func (wp *WorkerPool) Enqueue(task Task) bool {
	select {
	case wp.tasksChan <- task:
		return true
	default:
		return false // queue is full
	}
}

// waiting for items to hit the tasksChan
func (wp *WorkerPool) runWorker(ctx context.Context, workerID int) {
	defer wp.wg.Done()
	for {
		select {
		// return if context dies
		case <-ctx.Done():
			return
		// when a task hits this channel
		case task, ok := <-wp.tasksChan:
			if !ok {
				return
			}
			// finally call the process
			wp.processTask(ctx, workerID, task)
		}
	}
}

// calls the function initially passed in
func (wp *WorkerPool) processTask(ctx context.Context, workerID int, task Task) {
	log.Printf("[Worker %d] Processing video %s", workerID, task.VideoID)

	err := wp.enrichFn(ctx, task.VideoID)
	if err != nil {
		log.Printf("[Worker %d] Failed to enrich video %s: %v", workerID, task.VideoID, err)
		return
	}

	log.Printf("[Worker %d] Successfully enriched video %s", workerID, task.VideoID)
}

// stop closes the channel and waits for all active workers to finish
func (wp *WorkerPool) Stop() {
	close(wp.tasksChan)
	// block execution until counter reaches zero
	wp.wg.Wait()
}
