package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPutThenGet(t *testing.T) {
	cache := NewCache[string, string](10, 1*time.Millisecond)
	previouslyExisted := cache.Put("greeting", "hello")
	require.False(t, previouslyExisted)

	time.Sleep(3 * time.Millisecond)

	value, present := cache.Get("greeting")
	require.True(t, present)
	require.Equal(t, "hello", *value)
}

func TestGetMissing(t *testing.T) {
	cache := NewCache[string, string](1, 1*time.Millisecond)
	value, present := cache.Get("greeting")
	require.False(t, present)
	require.Nil(t, value)
}

func TestEviction_JustWrites(t *testing.T) {
	cache := NewCache[string, string](10, 1*time.Millisecond)

	for i := 0; i < 10; i++ {
		cache.Put(fmt.Sprintf("entry-%d", i), "hello")
	}

	time.Sleep(3 * time.Millisecond)

	_, present0 := cache.Get("entry-0")
	require.True(t, present0)

	_, present10 := cache.Get("entry-9")
	require.True(t, present10)

	cache.Put("entry-10", "hello")

	time.Sleep(3 * time.Millisecond)

	_, present1 := cache.Get("entry-1")
	require.False(t, present1)
}
