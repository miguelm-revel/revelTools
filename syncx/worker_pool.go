package syncx

import "sync"

// WorkerPool manages a fixed number of worker goroutines that process submitted jobs concurrently.
// Errors returned by jobs are collected and returned by Wait.
type WorkerPool struct {
	jobs chan func() error
	wg   sync.WaitGroup

	mu   sync.Mutex
	errs []error
}

// NewWorkerPool creates a WorkerPool with workers goroutines and a job channel of queueSize capacity.
func NewWorkerPool(workers, queueSize int) *WorkerPool {
	p := &WorkerPool{
		jobs: make(chan func() error, queueSize),
	}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				if err := job(); err != nil {
					p.mu.Lock()
					p.errs = append(p.errs, err)
					p.mu.Unlock()
				}
			}
		}()
	}
	return p
}

// Submit sends a job to the pool. Blocks if the job queue is full.
func (p *WorkerPool) Submit(job func() error) {
	p.jobs <- job
}

// Close closes the job channel, signalling workers to stop after draining remaining jobs.
func (p *WorkerPool) Close() {
	close(p.jobs)
}

// Wait blocks until all workers have finished and returns any errors collected from jobs.
func (p *WorkerPool) Wait() []error {
	p.wg.Wait()
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]error, len(p.errs))
	copy(out, p.errs)
	return out
}
