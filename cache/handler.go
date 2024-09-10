package memory

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	CacheDefaultCleanupTime    = 1 * time.Minute
	CacheDefaultExpirationTime = 1 * time.Minute
)

func (c *Cache) Get(key string, result interface{}) error {
	if !c.cacheIsActive {
		return fmt.Errorf("cache not active")
	}

	value, ok := c.cache.Load(key)
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	item, ok := value.(*cacheItem)
	if !ok || time.Now().After(item.expiration) {
		c.cache.Delete(key)
		return fmt.Errorf("key %s expired or invalid type", key)
	}

	return json.Unmarshal(item.data, result)
}

func (c *Cache) Set(key string, value interface{}) {
	if !c.cacheIsActive {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	item := &cacheItem{
		data:       data,
		expiration: time.Now().Add(CacheDefaultExpirationTime),
	}
	c.cache.Store(key, item)
}

func (c *Cache) Delete(key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *Cache) Clear() error {
	c.cache.Range(func(key, _ interface{}) bool {
		c.cache.Delete(key)
		return true
	})

	return nil
}

func (c *Cache) CleanupExpiredItems() {
	c.cache.Range(func(key, value interface{}) bool {
		item, ok := value.(*cacheItem)
		if !ok {
			c.cache.Delete(key)
			return true
		}

		if time.Now().After(item.expiration) {
			c.cache.Delete(key)
		}

		return true
	})
}

func (c *Cache) StartCleanupRoutine() {
	fmt.Println("[Memory-Cache] Cache cleanup routine started")

	ticker := time.NewTicker(CacheDefaultCleanupTime)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("[Memory-Cache] Cache cleanup executed")
		c.CleanupExpiredItems()
	}
}
