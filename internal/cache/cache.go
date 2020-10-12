package cache

import (
	"sync"
)

type Cache interface {
	Set(int64, string)
	Get(int64) (string, bool)
}

type InMemoryCache struct {
	sync.RWMutex
	items map[int64]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{items: make(map[int64]string)}
}

func (c *InMemoryCache) Set(id int64, value string) {
	c.Lock()
	c.items[id] = value
	c.Unlock()
}

func (c *InMemoryCache) Get(id int64) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.items[id]
	if !ok {
		return "", false
	}

	return value, true
}
