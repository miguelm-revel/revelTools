package sync

// Semaphore is a counting semaphore implemented using a buffered channel.
//
// The capacity of the channel defines the maximum number of goroutines
// that can hold the semaphore concurrently.
type Semaphore chan struct{}

// Lock acquires the semaphore, blocking if the maximum number of
// concurrent holders has been reached.
func (s Semaphore) Lock() {
	s <- struct{}{}
}

// Unlock releases the semaphore, allowing another waiting goroutine
// to acquire it.
func (s Semaphore) Unlock() {
	<-s
}

// NewSemaphore creates a new Semaphore that allows up to n concurrent holders.
func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}
