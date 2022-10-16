package cache

import (
	"crypto/sha256"
	"sync"
)

// This package provides a very simple cache. It's designed to hide the values of the keys because
// they will be used for storing authentication information, so keys are hashed before being used.
//
// 	c := Cache[int]()
// 	k := c.Key("secret number")
// 	c.Put(k, 42)
// 	...
// 	if v, ok := c.Get(k); ok {
//		...
// 	}

type Key [32]byte

type Entry[Value any] struct {
	value *Value
}

type Cache[Value any] struct {
	entries *sync.Map
}

func New[Value any]() *Cache[Value] {
	return &Cache[Value]{
		entries: &sync.Map{},
	}
}

func (c *Cache[V]) Key(k string) Key {
	return sha256.Sum256([]byte(k))
}

func (c *Cache[Value]) Get(k Key) (*Value, bool) {
	if value, ok := c.entries.Load(k); ok {
		if entry, ok := value.(Entry[Value]); ok {
			return entry.value, true
		}
	}
	return nil, false
}

func (c *Cache[Value]) Put(k Key, v *Value) {
	c.entries.Store(k, Entry[Value]{
		value: v,
	})
}
