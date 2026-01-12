package sync

type Semaphore chan struct{}

func (s Semaphore) Lock() {
	s <- struct{}{}
}

func (s Semaphore) Unlock() {
	<-s
}

func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}
