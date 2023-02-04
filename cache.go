package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type KeyExistsError struct {
	Key any
}

func (e *KeyExistsError) Error() string {
	return fmt.Sprintf("key %v already exists", e.Key)
}

type Cache struct {
	durationTimeEvict   time.Duration
	durationCheckTicker time.Duration
	store               map[any]*cacheValue
	mu                  *sync.RWMutex
}

type cacheValue struct {
	value    any
	lastUsed time.Time
}

func New(timeCheckNewTicker time.Duration, timeRecordEvict time.Duration) *Cache {
	return &Cache{
		mu:                  &sync.RWMutex{},
		store:               make(map[any]*cacheValue),
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

func (c *Cache) Add(key any, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.store[key]; ok {
		return &KeyExistsError{Key: key}
	}

	c.store[key] = &cacheValue{value: value, lastUsed: time.Now()}

	return nil
}

func (c *Cache) Get(key any) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, foundKey := c.store[key]

	if foundKey {
		if time.Since(val.lastUsed) >= c.durationTimeEvict {
			return nil, false
		}

		val.lastUsed = time.Now()

		return val.value, foundKey
	}

	return nil, false
}

func (c *Cache) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

func (c *Cache) Evict() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, val := range c.store {
		if time.Since(val.lastUsed) >= c.durationTimeEvict {
			delete(c.store, key)
		}
	}
}
