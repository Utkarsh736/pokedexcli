package pokecache

import (
	"sync"
	"time"
)

// Cache holds a map of cache entries protected by a mutex
type Cache struct {
	entries  map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache creates a new cache with a reaping interval
func NewCache(interval time.Duration) Cache {
	c := Cache{
		entries:  make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
		interval: interval,
	}
	
	// Start the reaping loop in a separate goroutine
	go c.reapLoop()
	
	return c
}

// Add inserts a new value into the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// reapLoop removes old entries every interval
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}

