package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type ICacheService interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type CacheService struct {
	client *memcache.Client
}

func (c *CacheService) Get(key string) (string, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}

func (c *CacheService) Set(key, value string) error {
	return c.client.Set(&memcache.Item{
		Key:   key,
		Value: []byte(value),
	})
}

func NewCacheService(port string) *CacheService {
	host := fmt.Sprintf("localost:%s", port)
	return &CacheService{client: memcache.New(host)}
}
