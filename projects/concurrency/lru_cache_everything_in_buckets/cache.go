package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	entryLimit int

	mu                     sync.Mutex
	entries                map[K]*cacheEntry[V]
	unsuccessfulReads      uint64
	evictedSuccessfulReads uint64
	evicted                uint64
	evictedNeverRead       uint64
}

type cacheEntry[V any] struct {
	value      V
	lastAccess time.Time
	reads      uint64
}

func NewCache[K comparable, V any](entryLimit int) *Cache[K, V] {
	return &Cache[K, V]{
		entryLimit: entryLimit,
		entries:    make(map[K]*cacheEntry[V]),
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether a value already existed in the cache for that key.
// If there was previously a value, it replaces that value with this one.
// Any Put counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.entries) == c.entryLimit {
		c.evict_locked()
	}

	entry, alreadyPresent := c.entries[key]
	if !alreadyPresent {
		entry = &cacheEntry[V]{
			value:      value,
			lastAccess: time.Now(),
		}
	} else {
		entry.value = value
		entry.lastAccess = time.Now()
	}

	c.entries[key] = entry

	return alreadyPresent
}

// evict_locked removes the oldest entry from the cache.
// We name this _locked to show that it assumes c.mu is held by the caller when this is called.
//
// This function is very expensive - we need to look through every element in the cache to decide whether it's the one to evict.
// This is O(n), which is pretty bad. Because we potentially perform an eviction for every write, this means that writing to the cache is O(n).
func (c *Cache[K, V]) evict_locked() {
	if len(c.entries) == 0 {
		return
	}
	isFirst := true
	var oldestKey K
	var oldestTimestamp time.Time
	for k, v := range c.entries {
		if isFirst || v.lastAccess.Before(oldestTimestamp) {
			oldestKey = k
			oldestTimestamp = v.lastAccess
			isFirst = false
		}
	}
	toEvict := c.entries[oldestKey]
	c.evictedSuccessfulReads += toEvict.reads
	if toEvict.reads == 0 {
		c.evictedNeverRead++
	}
	c.evicted++
	delete(c.entries, oldestKey)
}

// Get returns the value assocated with the passed key, and a boolean to indicate whether a value was known or not. If not, nil is returned as the value.
// Any Get counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, present := c.entries[key]
	if present {
		entry.lastAccess = time.Now()
		entry.reads++
		return &entry.value, present
	} else {
		c.unsuccessfulReads++
		return nil, false
	}
}

func (c *Cache[K, V]) Stats() CacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	writtenNeverRead := c.evictedNeverRead

	var readsCurrentValues uint64

	for _, entry := range c.entries {
		readsCurrentValues += entry.reads
		if entry.reads == 0 {
			writtenNeverRead++
		}
	}

	currentSize := uint64(len(c.entries))

	return CacheStats{
		sucessfulReadsAllTime:    c.evictedSuccessfulReads + readsCurrentValues,
		unsuccessfulReadsAllTime: c.unsuccessfulReads,
		writtenNeverReadAllTime:  writtenNeverRead,
		writesAllTime:            c.evicted + currentSize,
		readsCurrentValues:       readsCurrentValues,
		currentSize:              currentSize,
	}
}

type CacheStats struct {
	sucessfulReadsAllTime    uint64
	unsuccessfulReadsAllTime uint64
	writtenNeverReadAllTime  uint64
	writesAllTime            uint64
	readsCurrentValues       uint64
	currentSize              uint64
}

func (c *CacheStats) HitRate() float64 {
	return float64(c.sucessfulReadsAllTime) / (float64(c.sucessfulReadsAllTime) + float64(c.unsuccessfulReadsAllTime))
}

func (c *CacheStats) WrittenNeverRead() uint64 {
	return c.writtenNeverReadAllTime
}

func (c *CacheStats) AverageReadCountForCurrentEntries() float64 {
	return float64(c.readsCurrentValues) / float64(c.currentSize)
}

func (c *CacheStats) TotalReads() uint64 {
	return c.sucessfulReadsAllTime + c.unsuccessfulReadsAllTime
}

func (c *CacheStats) TotalSuccessfulReads() uint64 {
	return c.sucessfulReadsAllTime
}

func (c *CacheStats) TotalWrites() uint64 {
	return c.writesAllTime
}
