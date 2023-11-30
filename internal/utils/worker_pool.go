package utils

import (
	"go.uber.org/zap"
	"sync"
)

type WorkerPool struct {
	workers   int
	taskQueue chan func() error
	errChanel chan error
	logger    *zap.Logger
	wg        sync.WaitGroup
}

func NewWorkerPool(workers int, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func() error),
		logger:    logger,
		errChanel: make(chan error),
		wg:        sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start() {
	go wp.logError()
	for i := 0; i < wp.workers; i++ {
		go wp.runWorker()
	}
}

func (wp *WorkerPool) runWorker() {
	for task := range wp.taskQueue {
		err := task()
		if err != nil {
			wp.errChanel <- err
		}
		wp.wg.Done()
	}
}

func (wp *WorkerPool) logError() {
	for err := range wp.errChanel {
		wp.logger.Error("worker error", zap.Error(err))
	}
	wp.wg.Wait() // wait for all workers to finish & and all errors to be logged
	close(wp.errChanel)
}

func (wp *WorkerPool) Submit(task func() error) {
	wp.wg.Add(1)
	wp.taskQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}
