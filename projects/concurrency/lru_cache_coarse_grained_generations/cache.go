package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// type Cache implements a roughly-LRU cache. It attempts to keep to a maximum of targetSize, but may contain more entries at points in time.
// When under size pressure, it garbage collects entries which haven't been read or written, with no strict eviction ordering guarantees.
type Cache[K comparable, V any] struct {
	targetSize uint64

	mu sync.RWMutex
	// Every time we Get/Put a value, we store which generation it was last accessed.
	// We have a garbage collection goroutine which will delete entries that haven't been recently accessed, if the cache is full.
	currentGeneration atomic.Uint64
	values            map[K]*valueAndGeneration[V]
}

func NewCache[K comparable, V any](targetSize uint64, garbageCollectionInterval time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		targetSize: targetSize,
		values:     make(map[K]*valueAndGeneration[V], targetSize),
	}

	go func() {
		ticker := time.Tick(garbageCollectionInterval)
		for range ticker {
			currentGeneration := cache.currentGeneration.Load()
			cache.currentGeneration.Add(1)

			// Accumulate a keysToDelete slice so that we can collect the keys to delete under a read lock rather than holding a write lock for the entire GC cycle.
			// This will use extra memory, and has a disadvantage that we may bump a generation from a Get but then still evict that value because we already decided to GC it.
			var keysToDelete []K
			cache.mu.RLock()
			// If we have free space, don't garbage collect at all. This will probably lead to very spiky evictions.
			if uint64(len(cache.values)) <= targetSize {
				cache.mu.RUnlock()
				continue
			}
			for k, v := range cache.values {
				// This is a _very_ coarse-grained eviction policy. As soon as our cache becomes full, we may evict lots of entries.
				// It may be more useful to treat different values of generation differently, e.g. always evict if v.generation < currentGeneration - 5, and only evict more recent entries if that didn't free up any space.
				if v.generation.Load() != currentGeneration {
					keysToDelete = append(keysToDelete, k)
				}
			}
			cache.mu.RUnlock()
			if len(keysToDelete) > 0 {
				cache.mu.Lock()
				for _, keyToDelete := range keysToDelete {
					delete(cache.values, keyToDelete)
				}
				cache.mu.Unlock()
			}
		}
	}()

	return cache
}

type valueAndGeneration[V any] struct {
	value      V
	generation atomic.Uint64
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	valueWrapper := &valueAndGeneration[V]{
		value: value,
	}
	valueWrapper.generation.Store(c.currentGeneration.Load())
	c.values[key] = valueWrapper
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	valueWrapper, ok := c.values[key]
	if !ok {
		return nil, false
	}
	valueWrapper.generation.Store(c.currentGeneration.Load())
	return &valueWrapper.value, true
}
