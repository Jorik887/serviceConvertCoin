package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mu        sync.RWMutex
	data      map[string]float64
	expiresAt map[string]time.Time
}

// NewCache создает новый кэш
func NewCache() *Cache {
	return &Cache{
		data:      make(map[string]float64),
		expiresAt: make(map[string]time.Time),
	}
}

// Set добавляет курс валюты в кэш с указанием времени жизни.
func (c *Cache) Set(currency string, rate float64, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[currency] = rate
	c.expiresAt[currency] = time.Now().Add(duration)
}

// Get получает курс валюты из кэша
func (c *Cache) Get(currency string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if expiration, exists := c.expiresAt[currency]; exists {
		if time.Now().After(expiration) {
			// Если срок действия истек, удаляем из кеша
			c.mu.Lock() // Блокируем для удаления
			defer c.mu.Unlock()
			delete(c.data, currency)
			delete(c.expiresAt, currency)
			return 0, false
		}
		rate, exists := c.data[currency]
		return rate, exists
	}
	return 0, false
}
