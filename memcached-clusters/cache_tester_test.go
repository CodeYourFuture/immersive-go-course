package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayFlag(t *testing.T) {
	t.Run("ArrayFlag_String", func(t *testing.T) {
		a := arrayFlag([]string{"a", "b", "c"})

		require.Equal(t, a.String(), "a,b,c")
	})

	t.Run("ArrayFlag_Set", func(t *testing.T) {
		var a arrayFlag
		require.NoError(t, a.Set("a,b,c"))

		require.Equal(t, a, arrayFlag([]string{"a", "b", "c"}))
	})
}

type MockCacheService struct {
	cacheStore map[string]string
}

func (m *MockCacheService) Get(key string) (string, error) {
	return m.cacheStore[key], nil
}

func (m *MockCacheService) Set(key, value string) error {
	m.cacheStore[key] = value
	return nil
}

func TestMain(t *testing.T) {

	t.Run("Test sharded", func(t *testing.T) {
		router := &MockCacheService{cacheStore: map[string]string{}}
		nodes := map[string]ICacheService{
			"11211": &MockCacheService{cacheStore: map[string]string{
				"foo": "bar",
			}},
			"11212": &MockCacheService{cacheStore: map[string]string{
				"foo1": "not bar",
			}},
		}

		require.Equal(t, checkIfSharded(router, nodes), SHARDED)
	})

	t.Run("Test replicated", func(t *testing.T) {

		router := &MockCacheService{cacheStore: map[string]string{}}
		nodes := map[string]ICacheService{
			"11211": &MockCacheService{cacheStore: map[string]string{
				"foo": "bar",
			}},
			"11212": &MockCacheService{cacheStore: map[string]string{
				"foo": "bar",
			}},
		}

		require.Equal(t, checkIfSharded(router, nodes), REPLICATED)

	})
}
