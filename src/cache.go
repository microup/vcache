package cache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	durationTimeEvict   time.Duration
	durationCheckTicker time.Duration
	data                map[any] *cacheValue
	mu                  *sync.RWMutex
}

type cacheValue struct {
	value    any
	lastUsed time.Time
}

func New(timeCheckNewTicker time.Duration, timeRecordEvict time.Duration) *Cache {
	return &Cache{
		mu:                  &sync.RWMutex{},
		data:                make(map[any] *cacheValue),
		durationCheckTicker: timeCheckNewTicker,
		durationTimeEvict:   timeRecordEvict,		
	}
}

func (c *Cache) StartEvict(ctx context.Context) {
	ticker := time.NewTicker(c.durationCheckTicker)

	go func() {
		for {
			select {
			case <-ticker.C:
				c.Evict()
			case <-ctx.Done():
				ticker.Stop()

				return
			}
		}
	}()
}

func (c *Cache) Add(key any, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheValue{value: value, lastUsed: time.Now()}
}

func (c *Cache) Get(key any) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, foundKey := c.data[key]

	if foundKey {
		val.lastUsed = time.Now()
		c.data[key] = val

		return val.value, foundKey
	}

	return "", foundKey
}

func (c *Cache) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) Evict() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, val := range c.data {
		if time.Since(val.lastUsed) >= c.durationTimeEvict {
			delete(c.data, key)
		}
	}
}
