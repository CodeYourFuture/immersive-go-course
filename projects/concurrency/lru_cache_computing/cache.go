package cache

import (
	"container/list"
	"sync"
)

type Cache[K comparable, V any] struct {
	entryLimit uint64

	computeChannel chan K

	mu              sync.Mutex
	computedEntries map[K]cacheEntry[K, V]
	pendingEntries  map[K]*channelList[K, V]
	// Front is most recently used, back is least recently used
	evictionList *list.List
}

// entryLimit and concurrentComputeLimit must both be non-zero.
// computer must never panic.
func NewCache[K comparable, V any](entryLimit uint64, concurrentComputeLimit uint64, computer func(K) V) *Cache[K, V] {
	computeChannel := make(chan K, concurrentComputeLimit)

	resultChannel := make(chan keyValuePair[K, V], concurrentComputeLimit)

	for i := 0; i < int(concurrentComputeLimit); i++ {
		go func() {
			for key := range computeChannel {
				value := computer(key)
				resultChannel <- keyValuePair[K, V]{
					key:   key,
					value: &value,
				}
			}
		}()
	}

	cache := &Cache[K, V]{
		entryLimit:     entryLimit,
		computeChannel: computeChannel,

		computedEntries: make(map[K]cacheEntry[K, V], entryLimit),
		pendingEntries:  make(map[K]*channelList[K, V]),
		evictionList:    list.New(),
	}

	go func() {
		for result := range resultChannel {
			cache.mu.Lock()
			pendingEntry := cache.pendingEntries[result.key]
			delete(cache.pendingEntries, result.key)

			if len(cache.computedEntries) == int(cache.entryLimit) {
				keyToEvict := cache.evictionList.Remove(cache.evictionList.Back()).(K)
				delete(cache.computedEntries, keyToEvict)
			}

			evictionListPointer := cache.evictionList.PushFront(result.key)

			cache.computedEntries[result.key] = cacheEntry[K, V]{
				evictionListPointer: evictionListPointer,
				value:               *result.value,
			}
			pendingEntry.mu.Lock()
			pendingEntry.value = result.value
			cache.mu.Unlock()
			for _, ch := range pendingEntry.channels {
				ch <- result
			}
			pendingEntry.mu.Unlock()
		}
	}()

	return cache
}

type cacheEntry[K any, V any] struct {
	evictionListPointer *list.Element
	value               V
}

type keyValuePair[K any, V any] struct {
	key   K
	value *V
}

type channelList[K any, V any] struct {
	mu       sync.Mutex
	channels []chan (keyValuePair[K, V])
	value    *V
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	computedEntry, isComputed := c.computedEntries[key]
	pendingEntry, isPending := c.pendingEntries[key]
	if isComputed {
		c.evictionList.MoveToFront(computedEntry.evictionListPointer)
		c.mu.Unlock()
		return computedEntry.value, true
	}
	if !isPending {
		pendingEntry = &channelList[K, V]{}
		c.pendingEntries[key] = pendingEntry
	}
	c.mu.Unlock()
	if !isPending {
		c.computeChannel <- key
	}

	pendingEntry.mu.Lock()
	// Maybe the value was computed but hasn't been transfered from pending to computed yet
	if pendingEntry.value != nil {
		pendingEntry.mu.Unlock()
		return *pendingEntry.value, isPending
	}
	channel := make(chan keyValuePair[K, V], 1)
	pendingEntry.channels = append(pendingEntry.channels, channel)
	pendingEntry.mu.Unlock()
	value := <-channel
	return *value.value, isPending
}

// Only exists for testing. Doesn't count as a usage for LRU purposes.
func (c *Cache[K, V]) has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.computedEntries[key]
	return ok
}
