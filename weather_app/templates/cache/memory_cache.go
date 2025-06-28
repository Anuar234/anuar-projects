package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data   interface{}
	Expiry time.Time
}

type MemoryCache struct {
	store map[string]CacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		store: make(map[string]CacheEntry),
		ttl:   ttl,
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, exists := c.store[key]
	if !exists || time.Now().After(entry.Expiry) {
		return nil, false
	}
	return entry.Data, true
}

func (c *MemoryCache) Set(key string, data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = CacheEntry{
		Data:   data,
		Expiry: time.Now().Add(c.ttl),
	}
}

func (c *MemoryCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.store {
		if time.Now().After(entry.Expiry) {
			delete(c.store, key)
		}
	}
}
