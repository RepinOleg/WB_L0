package cache

import "sync"

type Cache struct {
	data map[uint64]string
	mu   sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		data: make(map[uint64]string),
	}
	return cache
}

func (c *Cache) SetOrder(id uint64, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = data
}

func (c *Cache) GetOrderByID(id uint64) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok := c.data[id]
	if !ok {
		return data, false
	}
	return data, true
}
