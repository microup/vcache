package cache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	durationTimeEvict   time.Duration
	durationCheckTicker time.Duration
	data                map[string] cacheValue
	mu                  *sync.Mutex
}

type cacheValue struct {
	value    interface{}
	lastUsed time.Time
}

func New(timeCheckNewTicker time.Duration, timeRecordEvict time.Duration) *Cache {
	return &Cache{
		mu:                  &sync.Mutex{},
		data:                make(map[string]cacheValue),
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

func (c *Cache) Add(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheValue{value: value, lastUsed: time.Now()}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, foundKey := c.data[key]

	if foundKey {
		val.lastUsed = time.Now()
		c.data[key] = val

		return val.value, foundKey
	}

	return "", foundKey
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
