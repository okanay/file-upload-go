package memory

import (
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	data       []byte
	expiration time.Time
}

type Cache struct {
	cache         sync.Map
	cacheIsActive bool
}

func Init(cacheIsActive bool) *Cache {
	fmt.Println("[Memory Cache] Initializing the memory cache")

	c := &Cache{
		cache:         sync.Map{},
		cacheIsActive: cacheIsActive,
	}

	if cacheIsActive {
		go c.StartCleanupRoutine()
	}

	return c
}
