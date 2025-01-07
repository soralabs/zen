package cache

import (
	"context"
	"sync/atomic"
	"time"
)

var (
	hits    int64
	misses  int64
	evicted int64
)

func New(config Config) *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{
		items:   make(map[CacheKey]CacheEntry),
		maxSize: config.MaxSize,
		ttl:     config.TTL,
		ctx:     ctx,
		cancel:  cancel,
	}

	// Start cleanup routine
	go c.cleanup(config.CleanupPeriod)
	return c
}

func (c *Cache) Set(key CacheKey, value interface{}) {
	c.Lock()
	defer c.Unlock()

	// Check if we need to evict
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	c.items[key] = CacheEntry{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key CacheKey) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	entry, exists := c.items[key]
	if !exists {
		atomic.AddInt64(&misses, 1)
		return nil, false
	}

	if time.Now().After(entry.Expiration) {
		atomic.AddInt64(&misses, 1)
		return nil, false
	}

	atomic.AddInt64(&hits, 1)
	return entry.Value, true
}

func (c *Cache) Delete(key CacheKey) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[CacheKey]CacheEntry)
}

func (c *Cache) GetStats() CacheStats {
	c.RLock()
	defer c.RUnlock()

	return CacheStats{
		Size:    len(c.items),
		Hits:    atomic.LoadInt64(&hits),
		Misses:  atomic.LoadInt64(&misses),
		Evicted: atomic.LoadInt64(&evicted),
	}
}

func (c *Cache) cleanup(period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Lock()
			now := time.Now()
			for key, entry := range c.items {
				if now.After(entry.Expiration) {
					delete(c.items, key)
					atomic.AddInt64(&evicted, 1)
				}
			}
			c.Unlock()
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Cache) evictOldest() {
	var oldestKey CacheKey
	var oldestTime time.Time

	for key, entry := range c.items {
		if oldestTime.IsZero() || entry.Expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Expiration
		}
	}

	if !oldestTime.IsZero() {
		delete(c.items, oldestKey)
		atomic.AddInt64(&evicted, 1)
	}
}

func (c *Cache) Close() {
	c.cancel()
}
