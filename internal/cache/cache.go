package cache

import (
	"github.com/RepinOleg/WB_L0/internal/models"
	"sync"
)

type Cache struct {
	data map[string]models.Order
	mu   sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		data: make(map[string]models.Order),
	}
	return cache
}

func (c *Cache) SetOrder(id string, data models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = data
}

func (c *Cache) GetOrderByID(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok := c.data[id]
	if !ok {
		return data, false
	}
	return data, true
}

func (c *Cache) SetAllOrders(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.data[order.OrderUID] = order
	}
}

func (c *Cache) GetAllOrders() []models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]models.Order, 0, len(c.data))
	for _, order := range c.data {
		res = append(res, order)
	}
	return res
}
