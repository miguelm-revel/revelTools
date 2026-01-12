package sync

import "sync"

type Barrier struct{ *sync.WaitGroup }

func (b *Barrier) Lock() {
	b.Add(1)
}

func (b *Barrier) Unlock() {
	b.Done()
	b.Wait()
}

func NewBarrier() *Barrier {
	return &Barrier{&sync.WaitGroup{}}
}
