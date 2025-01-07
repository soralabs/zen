package cache

import (
	"context"
	"sync"
	"time"
)

type CacheKey string

type CacheEntry struct {
	Value      interface{}
	Expiration time.Time
}

type Cache struct {
	sync.RWMutex
	items   map[CacheKey]CacheEntry
	maxSize int
	ttl     time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

type Config struct {
	MaxSize       int
	TTL           time.Duration
	CleanupPeriod time.Duration
}

type CacheStats struct {
	Size    int
	Hits    int64
	Misses  int64
	Evicted int64
}
