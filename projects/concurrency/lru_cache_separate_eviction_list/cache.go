package cache

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type Cache[K comparable, V any] struct {
	entryLimit int

	writes           atomic.Uint64
	hits             atomic.Uint64
	misses           atomic.Uint64
	writtenNeverRead atomic.Int64

	mu     sync.Mutex
	values map[K]*list.Element
	// Front is most recent element, Back is next to be evicted.
	// We use a *list.List here because moving elements within a *list.List, as well as adding a new element at one end, or finding the element at one end, is O(1), so all of the following are cheap:
	// * Adding a new element to the eviction list.
	// * Finding which element should be evicted next.
	// * Moving an element from anywhere in the list to one end of it .
	// The values inside each list element all have type *keyAndValueContainer[K, V].
	// Ideally this would be a generic type, but it pre-dates generics in the language,
	// so we need to ourselves track what type we expect to put in here, and we need to assert the type when we get values out.
	evictionList *list.List
}

type keyAndValueContainer[K any, V any] struct {
	key       K
	value     V
	readCount atomic.Uint64
}

func NewCache[K comparable, V any](entryLimit int) *Cache[K, V] {
	evictionList := list.New()
	return &Cache[K, V]{
		entryLimit: entryLimit,
		// Pre-allocate the whole value map - this optimises for assuming our cache will be pretty full - if we expected it to be mostly empty, we may not pre-allocate here.
		values:       make(map[K]*list.Element, entryLimit),
		evictionList: evictionList,
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether a value already existed in the cache for that key.
// If there was previously a value, it replaces that value with this one.
// Any Put counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, alreadyKnown := c.values[key]
	if !alreadyKnown && len(c.values) == c.entryLimit {
		c.evict_locked()
	}
	keyAndValue := &keyAndValueContainer[K, V]{
		key:   key,
		value: value,
	}
	c.writes.Add(1)
	c.writtenNeverRead.Add(1)
	c.values[key] = c.evictionList.PushFront(keyAndValue)
	return alreadyKnown
}

// refresh_locked moves a particular key to be the last element to be evicted if it's known (or returns nil if not).
// We name this _locked to show that it assumes c.mu is held by the caller when this is called.
func (c *Cache[K, V]) refresh_locked(key K) *keyAndValueContainer[K, V] {
	element, known := c.values[key]
	if !known {
		return nil
	}
	keyAndValue := c.evictionList.Remove(element)
	c.evictionList.PushFront(keyAndValue)
	return keyAndValue.(*keyAndValueContainer[K, V])
}

// evict_locked removes the oldest entry from the cache.
// We name this _locked to show that it assumes c.mu is held by the caller when this is called.
func (c *Cache[K, V]) evict_locked() {
	delete(c.values, c.evictionList.Remove(c.evictionList.Back()).(*keyAndValueContainer[K, V]).key)
}

// Get returns the value assocated with the passed key, and a boolean to indicate whether a value was known or not. If not, nil is returned as the value.
// Any Get counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if keyAndValue := c.refresh_locked(key); keyAndValue != nil {
		c.hits.Add(1)
		if keyAndValue.readCount.Add(1) == 1 {
			// If this is the first read of the value, this value has moved from being "written never read", to now having been read.
			c.writtenNeverRead.Add(-1)
		}
		return &keyAndValue.value, true
	}
	c.misses.Add(1)
	return nil, false
}

func (c *Cache[K, V]) Stats() CacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	var readsCurrentValues uint64
	value := c.evictionList.Front()
	for value != nil {
		readsCurrentValues += value.Value.(*keyAndValueContainer[K, V]).readCount.Load()
		value = value.Next()
	}

	return CacheStats{
		sucessfulReadsAllTime:    c.hits.Load(),
		unsuccessfulReadsAllTime: c.misses.Load(),
		writtenNeverReadAllTime:  uint64(c.writtenNeverRead.Load()),
		writesAllTime:            c.writes.Load(),
		readsCurrentValues:       readsCurrentValues,
		currentSize:              uint64(len(c.values)),
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
