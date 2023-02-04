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
	durationRecordEvict    int64
	durationCheckNewTicker time.Duration
	store                  map[any]*cacheValue
	mu                     *sync.RWMutex
}

type cacheValue struct {
	value          any
	expirationTime int64
}

func New(timeCheckNewTicker time.Duration, timeRecordEvict time.Duration) *Cache {
	return &Cache{
		mu:                     &sync.RWMutex{},
		store:                  make(map[any]*cacheValue),
		durationCheckNewTicker: timeCheckNewTicker,
		durationRecordEvict:    int64(timeRecordEvict.Seconds()),
	}
}

func (c *Cache) StartEvict(ctx context.Context) {
	ticker := time.NewTicker(c.durationCheckNewTicker)

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

	if _, ok := c.store[key]; ok {
		c.mu.Unlock()

		return &KeyExistsError{Key: key}
	}

	c.store[key] = &cacheValue{
		value:          value,
		expirationTime: time.Now().Unix()  + c.durationRecordEvict,
	}

	c.mu.Unlock()

	return nil
}

func (c *Cache) Get(key any) (any, bool) {
	c.mu.RLock()

	val, foundKey := c.store[key]

	if foundKey {
		if time.Now().Unix() > val.expirationTime {
			c.mu.RUnlock()

			return nil, false
		}

		val.expirationTime = time.Now().Unix()

		c.mu.RUnlock()

		return val.value, foundKey
	}

	c.mu.RUnlock()

	return nil, false
}

func (c *Cache) Delete(key any) {
	c.mu.Lock()

	delete(c.store, key)
	c.mu.Unlock()
}

func (c *Cache) Evict() {
	c.mu.Lock()

	for key, val := range c.store {
		if time.Now().Unix() > val.expirationTime {
			delete(c.store, key)
		}
	}

	c.mu.Unlock()
}
