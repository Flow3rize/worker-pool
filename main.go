package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Flow3rize/worker-pool/worker"
)

func main() {
	pool := worker.NewWorkerPool(10)

	var resultsWG sync.WaitGroup
	resultsWG.Add(1)
	go func() {
		defer resultsWG.Done()
		for res := range pool.Results {
			fmt.Printf("Результат: %s\n", res)
		}
		fmt.Println("Канал результатов закрыт")
	}()

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
