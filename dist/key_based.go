package dist

import (
	"hash/fnv"
	"sync"
)

// KeyBasedPool distributes items deterministically by key using FNV-1a hashing.
// The same key always maps to the same instance. Safe for concurrent use.
type KeyBasedPool[T any] struct {
	pool  []T
	mutex sync.RWMutex
}

// Pick returns the instance deterministically assigned to key via FNV-1a hash.
func (k *KeyBasedPool[T]) Pick(key string) T {
	h := fnv.New32a()
	h.Write([]byte(key))
	k.mutex.RLock()
	defer k.mutex.RUnlock()
	idx := h.Sum32() % uint32(len(k.pool))
	return k.pool[idx]
}

// NewKeyBasedPool creates a KeyBasedPool backed by items.
func NewKeyBasedPool[T any](items []T) *KeyBasedPool[T] {
	return &KeyBasedPool[T]{pool: items}
}

// KeyBasedPoolFromConstructor builds a KeyBasedPool of poolSize instances using constructor.
func KeyBasedPoolFromConstructor[T any](constructor func() T, poolSize int) *KeyBasedPool[T] {
	pool := make([]T, poolSize)
	for i := 0; i < poolSize; i++ {
		pool[i] = constructor()
	}
	return &KeyBasedPool[T]{pool: pool}
}
