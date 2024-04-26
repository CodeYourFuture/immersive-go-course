package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPutThenGet(t *testing.T) {
	gcTicker := make(chan time.Time)
	cache := NewCache[string, string](10, gcTicker)
	previouslyExisted := cache.Put("greeting", "hello")
	require.False(t, previouslyExisted)

	// Write to the channel twice twice, because we know that once the second write has sent the first one must be done processing.
	gcTicker <- time.Now()
	gcTicker <- time.Now()

	value, present := cache.Get("greeting")
	require.True(t, present)
	require.Equal(t, "hello", *value)
}

func TestGetMissing(t *testing.T) {
	gcTicker := make(chan time.Time)
	cache := NewCache[string, string](1, gcTicker)
	value, present := cache.Get("greeting")
	require.False(t, present)
	require.Nil(t, value)
}

func TestEviction_JustWrites(t *testing.T) {
	gcTicker := make(chan time.Time)
	cache := NewCache[string, string](10, gcTicker)

	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("entry-%d", i), "hello")
	}

	gcTicker <- time.Now()
	gcTicker <- time.Now()

	_, present0 := cache.Get("entry-0")
	require.True(t, present0)

	_, present10 := cache.Get("entry-9")
	require.True(t, present10)

	cache.Put("entry-10", "hello")

	gcTicker <- time.Now()
	gcTicker <- time.Now()

	presentCount := 0
	for key := 0; key < 11; key++ {
		got, present := cache.Get(fmt.Sprintf("entry-%d", key))
		if present {
			presentCount++
			require.Equal(t, "hello", *got)
		}
	}
	require.Equal(t, 10, presentCount)

	// entries 0, 9, and 10 were accessed a generation after the others, so should be present.
	_, present0AfterGC := cache.Get("entry-0")
	require.True(t, present0AfterGC)

	_, present9AfterGC := cache.Get("entry-9")
	require.True(t, present9AfterGC)

	_, present10AfterGC := cache.Get("entry-10")
	require.True(t, present10AfterGC)
}

func TestConcurrentWrites(t *testing.T) {
	gcTicker := make(chan time.Time)
	cache := NewCache[int, string](1, gcTicker)

	var wg sync.WaitGroup

	for iteration := 0; iteration < 100000; iteration++ {
		wg.Add(1)
		go func() {
			for key := 0; key < 3; key++ {
				cache.Put(key, fmt.Sprintf("entry-%d", key))
			}
			wg.Done()
		}()
	}

	wg.Wait()

	// No gc tick has happened, so all three keys should be present.
	got0, present0 := cache.Get(0)
	require.True(t, present0)
	require.Equal(t, "entry-0", *got0)

	got1, present1 := cache.Get(1)
	require.True(t, present1)
	require.Equal(t, "entry-1", *got1)

	got2, present2 := cache.Get(2)
	require.True(t, present2)
	require.Equal(t, "entry-2", *got2)

	gcTicker <- time.Now()
	gcTicker <- time.Now()

	presentCount := 0
	for key := 0; key < 3; key++ {
		got, present := cache.Get(key)
		if present {
			presentCount++
			require.Equal(t, fmt.Sprintf("entry-%d", key), *got)
		}
	}
	require.Equal(t, 1, presentCount)
}
