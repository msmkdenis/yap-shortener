package worker_pool

import (
	"go.uber.org/zap"
)

// WorkerPool represents worker pool
type WorkerPool struct {
	workers   int
	taskQueue chan func() error
	logger    *zap.Logger
}

// NewWorkerPool returns a new instance of WorkerPool
func NewWorkerPool(workers int, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func() error),
		logger:    logger,
	}
}

// Start starts async workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.runWorker()
	}
}

func (wp *WorkerPool) runWorker() {
	for task := range wp.taskQueue {
		err := task()
		if err != nil {
			wp.logger.Error("worker error", zap.Error(err))
		}
	}
}

// Submit submits task to worker pool
func (wp *WorkerPool) Submit(task func() error) {
	wp.taskQueue <- task
}

// Stop stops worker pool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}
