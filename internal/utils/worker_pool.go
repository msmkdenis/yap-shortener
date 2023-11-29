package utils

import (
	"sync"

	"go.uber.org/zap"
)

type WorkerPool struct {
	workers   int
	taskQueue chan func() error
	errChanel chan error
	wg        sync.WaitGroup
	logger    *zap.Logger
}

func NewWorkerPool(workers int, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func() error),
		logger:    logger,
		errChanel: make(chan error),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.runWorker()
	}
	wp.logError()
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
	for i := 0; i < wp.workers; i++ {
		go func() {
			for err := range wp.errChanel {
				wp.logger.Error("WorkPoolTaskError", zap.Error(err))
			}
		}()
	}
}

func (wp *WorkerPool) Submit(task func() error) {
	wp.wg.Add(1)
	wp.taskQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	close(wp.errChanel)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
