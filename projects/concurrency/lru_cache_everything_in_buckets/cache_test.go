package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutThenGet(t *testing.T) {
	cache := NewCache[string, string](10)
	previouslyExisted := cache.Put("greeting", "hello")
	require.False(t, previouslyExisted)

	value, present := cache.Get("greeting")
	require.True(t, present)
	require.Equal(t, "hello", *value)
}

func TestGetMissing(t *testing.T) {
	cache := NewCache[string, string](10)
	value, present := cache.Get("greeting")
	require.False(t, present)
	require.Nil(t, value)
}

func TestPutThenOverwriteSameValue(t *testing.T) {
	cache := NewCache[string, string](10)
	previouslyExisted1 := cache.Put("greeting", "hello")
	require.False(t, previouslyExisted1)

	previouslyExisted2 := cache.Put("greeting", "hello")
	require.True(t, previouslyExisted2)

	value, present := cache.Get("greeting")
	require.True(t, present)
	require.Equal(t, "hello", *value)
}

func TestPutThenOverwriteDifferentValue(t *testing.T) {
	cache := NewCache[string, string](10)
	previouslyExisted1 := cache.Put("greeting", "hello")
	require.False(t, previouslyExisted1)

	previouslyExisted2 := cache.Put("greeting", "howdy")
	require.True(t, previouslyExisted2)

	value, present := cache.Get("greeting")
	require.True(t, present)
	require.Equal(t, "howdy", *value)
}

func TestEviction_JustWrites(t *testing.T) {
	cache := NewCache[string, string](10)

	for i := 0; i <= 10; i++ {
		cache.Put(fmt.Sprintf("entry-%d", i), "hello")
	}
	_, present0 := cache.Get("entry-0")
	require.False(t, present0)

	_, present10 := cache.Get("entry-10")
	require.True(t, present10)
}

func TestEviction_ReadsAndWrites(t *testing.T) {
	cache := NewCache[string, string](10)

	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("entry-%d", i), "hello")
	}
	_, present0 := cache.Get("entry-0")
	require.True(t, present0)

	cache.Put("entry-10", "hello")

	_, present0 = cache.Get("entry-0")
	require.True(t, present0)

	_, present1 := cache.Get("entry-1")
	require.False(t, present1)

	_, present10 := cache.Get("entry-10")
	require.True(t, present10)
}

func TestConcurrentPuts(t *testing.T) {
	cache := NewCache[string, string](10)

	var startWaitGroup sync.WaitGroup
	var endWaitGroup sync.WaitGroup

	for i := 0; i < 1000; i++ {
		startWaitGroup.Add(1)
		endWaitGroup.Add(1)
		go func(i int) {
			startWaitGroup.Wait()
			cache.Put(fmt.Sprintf("entry-%d", i), "hello")
			endWaitGroup.Done()
		}(i)
	}
	startWaitGroup.Add(-1000)
	endWaitGroup.Wait()

	sawEntries := 0
	for i := 0; i < 1000; i++ {
		_, saw := cache.Get(fmt.Sprintf("entry-%d", i))
		if saw {
			sawEntries++
		}
	}
	require.Equal(t, 10, sawEntries)
}

func TestStats(t *testing.T) {
	cache := NewCache[string, string](10)

	var startWaitGroup sync.WaitGroup
	var endWaitGroup sync.WaitGroup

	for i := 0; i < 1000; i++ {
		startWaitGroup.Add(1)
		endWaitGroup.Add(1)
		go func(i int) {
			startWaitGroup.Wait()
			cache.Put(fmt.Sprintf("entry-%d", i), "hello")
			endWaitGroup.Done()
		}(i)
	}
	startWaitGroup.Add(-1000)
	endWaitGroup.Wait()

	var readWaitGroup sync.WaitGroup
	var sawEntries atomic.Uint64
	for i := 0; i < 1000; i++ {
		readWaitGroup.Add(1)
		go func(i int) {
			_, saw := cache.Get(fmt.Sprintf("entry-%d", i))
			if saw {
				sawEntries.Add(1)
			}
			readWaitGroup.Done()
		}(i)
	}
	readWaitGroup.Wait()
	require.Equal(t, uint64(10), sawEntries.Load())

	stats := cache.Stats()
	require.Equal(t, 0.01, stats.HitRate())
	require.Equal(t, uint64(990), stats.WrittenNeverRead())
	require.Equal(t, float64(1), stats.AverageReadCountForCurrentEntries())
	require.Equal(t, uint64(1000), stats.TotalReads())
	require.Equal(t, uint64(10), stats.TotalSuccessfulReads())
	require.Equal(t, uint64(1000), stats.TotalWrites())
}
