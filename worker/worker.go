package worker

import (
	"sync"
	"time"
)

type Worker struct {
	id       int
	jobs     <-chan string
	results  chan<- string
	isClosed chan bool
	wg       *sync.WaitGroup
}

type WorkerPool struct {
	jobs    chan string
	results chan string
	workers map[int]*Worker
	wg      sync.WaitGroup
	nextID  int
	mu      sync.Mutex
}

func NewWorker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:       id,
		jobs:     jobs,
		results:  results,
		isClosed: make(chan bool),
		wg:       wg,
	}
}

func NewWorkerPool(bufSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan string, bufSize),
		results: make(chan string, bufSize),
		workers: make(map[int]*Worker),
		nextID:  1,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case job, ok := <-w.jobs:
				if !ok {
					InfoLogger.Printf("Worker %d: channel closed, finishing work\n", w.id)
					return
				}
				InfoLogger.Printf("Worker %d processed: %s\n", w.id, job)
				time.Sleep(500 * time.Millisecond)
				w.results <- job
			case <-w.isClosed:
				InfoLogger.Printf("Worker %d received a stop signal\n", w.id)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.isClosed)
}

func (wp *WorkerPool) AddWorker() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	id := wp.nextID
	wp.nextID++

	worker := NewWorker(id, wp.jobs, wp.results, &wp.wg)
	wp.workers[id] = worker
	worker.Start()
	InfoLogger.Printf("Worker %d added", id)
	return id
}

func (wp *WorkerPool) RemoveWorker(id int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	worker, exists := wp.workers[id]
	if !exists {
		ErrorLogger.Printf("Worker %d not found\n", id)
		return
	}
	worker.Stop()
	delete(wp.workers, id)
	InfoLogger.Printf("Worker %d deleted", id)
}

func (wp *WorkerPool) Submit(job string) {
	wp.jobs <- job
}

func (wp *WorkerPool) Shutdown() {
	close(wp.jobs)
	wp.wg.Wait()
	InfoLogger.Printf("The pool has completed its work.")
}
