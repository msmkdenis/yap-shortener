package utils

import (
	"go.uber.org/zap"
)

type WorkerPool struct {
	workers   int
	taskQueue chan func() error
	logger    *zap.Logger
}

func NewWorkerPool(workers int, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func() error),
		logger:    logger,
	}
}

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

func (wp *WorkerPool) Submit(task func() error) {
	wp.taskQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}
