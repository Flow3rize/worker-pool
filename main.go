package main

import (
	"fmt"
	"time"

	"github.com/Flow3rize/worker-pool/worker"
)

func main() {
	pool := worker.NewWorkerPool(10)

	w1 := pool.AddWorker()
	w2 := pool.AddWorker()

	for i := 1; i <= 7; i++ {
		pool.Submit(fmt.Sprintf("Job %d", i))
	}

	time.Sleep(2 * time.Second)

	w3 := pool.AddWorker()

	for i := 8; i <= 14; i++ {
		pool.Submit(fmt.Sprintf("Job %d", i))
	}

	time.Sleep(2 * time.Second)

	pool.RemoveWorker(w2)

	for i := 15; i <= 21; i++ {
		pool.Submit(fmt.Sprintf("Job %d", i))
	}

	time.Sleep(3 * time.Second)

	pool.RemoveWorker(w1)
	pool.RemoveWorker(w3)

	pool.Shutdown()
}
