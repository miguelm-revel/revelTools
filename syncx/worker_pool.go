package syncx

import "sync"

type WorkerPool struct {
	jobs chan func() error
	wg   sync.WaitGroup

	mu   sync.Mutex
	errs []error
}

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

func (p *WorkerPool) Submit(job func() error) {
	p.jobs <- job
}

func (p *WorkerPool) Close() {
	close(p.jobs)
}

func (p *WorkerPool) Wait() []error {
	p.wg.Wait()
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]error, len(p.errs))
	copy(out, p.errs)
	return out
}
