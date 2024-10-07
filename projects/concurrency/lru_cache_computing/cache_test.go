package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetTwice(t *testing.T) {
	cache := NewCache[string, string](10, 1, func(string) string {
		time.Sleep(10 * time.Millisecond)
		return "hello"
	})

	timeBefore := time.Now()

	value0, previouslyExisted0 := cache.Get("greeting")
	require.False(t, previouslyExisted0)
	require.Equal(t, "hello", value0)

	value1, previouslyExisted1 := cache.Get("greeting")
	require.True(t, previouslyExisted1)
	require.Equal(t, "hello", value1)

	elapsedTime := time.Since(timeBefore)

	// Should have only been computed once
	require.Less(t, elapsedTime, 20*time.Millisecond)
}

func TestConcurrencyLimit(t *testing.T) {
	cache := NewCache[string, string](10, 1, func(string) string {
		time.Sleep(10 * time.Millisecond)
		return "hello"
	})

	var wg sync.WaitGroup
	wg.Add(2)

	timeBefore := time.Now()

	go func() {
		value0, previouslyExisted0 := cache.Get("greeting0")
		require.False(t, previouslyExisted0)
		require.Equal(t, "hello", value0)
		wg.Done()
	}()

	go func() {
		value1, previouslyExisted1 := cache.Get("greeting1")
		require.False(t, previouslyExisted1)
		require.Equal(t, "hello", value1)
		wg.Done()
	}()

	wg.Wait()

	elapsedTime := time.Since(timeBefore)

	require.Greater(t, elapsedTime, 19*time.Millisecond)
}

func TestGetTwoDifferentValuesInParallel(t *testing.T) {
	cache := NewCache[string, string](10, 2, func(string) string {
		time.Sleep(10 * time.Millisecond)
		return "hello"
	})

	var wg sync.WaitGroup
	wg.Add(2)

	timeBefore := time.Now()

	go func() {
		value0, previouslyExisted0 := cache.Get("greeting0")
		require.False(t, previouslyExisted0)
		require.Equal(t, "hello", value0)
		wg.Done()
	}()

	go func() {
		value1, previouslyExisted1 := cache.Get("greeting1")
		require.False(t, previouslyExisted1)
		require.Equal(t, "hello", value1)
		wg.Done()
	}()

	wg.Wait()

	elapsedTime := time.Since(timeBefore)

	// Should have only been computed once
	require.Less(t, elapsedTime, 20*time.Millisecond)
}

func TestEvictsInOrder(t *testing.T) {
	cache := NewCache[string, string](2, 2, func(string) string {
		return "hello"
	})

	value0, previouslyExisted0 := cache.Get("greeting0")
	require.False(t, previouslyExisted0)
	require.Equal(t, "hello", value0)

	value1, previouslyExisted1 := cache.Get("greeting1")
	require.False(t, previouslyExisted1)
	require.Equal(t, "hello", value1)

	require.True(t, cache.has("greeting0"))
	require.True(t, cache.has("greeting1"))

	value2, previespreviouslyExisted2 := cache.Get("greeting2")
	require.False(t, previespreviouslyExisted2)
	require.Equal(t, "hello", value2)

	require.False(t, cache.has("greeting0"))
	require.True(t, cache.has("greeting1"))
	require.True(t, cache.has("greeting2"))
}
