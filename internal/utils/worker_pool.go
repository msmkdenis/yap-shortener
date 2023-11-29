package utils

import (
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

type WorkerPool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
	logger    *zap.Logger
}

func NewWorkerPool(logger *zap.Logger) *WorkerPool {
	workersEnv := os.Getenv("WORKERS")
	workers, err := strconv.Atoi(workersEnv)
	if err != nil {
		workers = 100
	}

	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func()),
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
		task()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1)
	wp.taskQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
