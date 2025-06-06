package worker

import (
	"testing"
	"time"
)

func TestAddAndRemoveWorker(t *testing.T) {
	pool := NewWorkerPool(5)
	id := pool.AddWorker()
	if _, ok := pool.workers[id]; !ok {
		t.Errorf("Worker %d was not added", id)
	}
	pool.RemoveWorker(id)
	if _, ok := pool.workers[id]; ok {
		t.Errorf("Worker %d was not removed", id)
	}
}

func TestSubmitJob(t *testing.T) {
	pool := NewWorkerPool(2)
	id := pool.AddWorker()
	defer pool.RemoveWorker(id)

	job := "test job"
	pool.Submit(job)

	select {
	case res := <-pool.Results:
		if res != job {
			t.Errorf("Expected result %q, got %q", job, res)
		}
	case <-time.After(time.Second):
		t.Error("Job was not processed in time")
	}
}
