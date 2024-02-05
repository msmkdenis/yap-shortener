package workerpool

import (
	"fmt"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

// ExampleNewWorkerPool shows how to use WorkerPool
func ExampleNewWorkerPool() {
	// Initialize logger
	// Errors produced by workers will be logged with Error level
	logger, _ := zap.NewProduction()

	// Create worker pool with 3 workers
	pool := NewWorkerPool(3, logger)

	// Start async workers
	pool.Start()

	// Stop worker pool
	defer pool.Stop()

	var counter atomic.Int64

	var wg sync.WaitGroup
	// Submit 10 CAS operations to worker pool
	for i := 0; i < 10; i++ {
		wg.Add(1)
		pool.Submit(func() error {
			defer wg.Done()
			counter.Add(1)
			return nil
		})
	}

	// Wait for all CAS operations
	wg.Wait()
	fmt.Println(counter.Load())
	// Output: 10
}
