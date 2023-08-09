package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrContextIsNil = errors.New("context is nil")

//KeyExistsError: a custom error type that is returned when an attempt is
// made to add a key to the cache that already exists.
type KeyExistsError struct {
	Key any
}

func (e *KeyExistsError) Error() string {
	return fmt.Sprintf("key %v already exists", e.Key)
}

//VCache: a struct that implements a simple in-memory key-value cache with eviction.
type VCache struct {
	durationRecordEvict    int64
	durationCheckNewTicker time.Duration
	store                  map[any]*cacheValue
	mu                     *sync.RWMutex
}

//cacheValue: a struct that contains a value and its expiration time.
type cacheValue struct {
	value          any
	expirationTime int64
}

//New: a function that creates and returns a new Cache instance.
func New(timeCheckNewTicker time.Duration, timeRecordEvict time.Duration) *VCache {
	return &VCache{
		mu:                     &sync.RWMutex{},
		store:                  make(map[any]*cacheValue),
		durationCheckNewTicker: timeCheckNewTicker,
		durationRecordEvict:    int64(timeRecordEvict.Seconds()),
	}
}

//StartEvict: a method that starts the eviction process in a separate goroutine.
// It stops when the context passed as an argument is done.
func (c *VCache) StartEvict(ctx context.Context) error {
	if ctx == nil {
		return ErrContextIsNil
	}

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

	return nil
}

//Add: a method that adds a new key-value pair to the cache. Returns KeyExistsError if the key already exists.
func (c *VCache) Add(key any, value any) error {
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

//Get: a method that retrieves a value from the cache by its key.
// Returns the value and a boolean indicating if the key was found.
func (c *VCache) Get(key any) (any, bool) {
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

//Delete: a method that deletes a key-value pair from the cache.
func (c *VCache) Delete(key any) {
	c.mu.Lock()

	delete(c.store, key)
	c.mu.Unlock()
}

//Evict: a method that evicts expired key-value pairs from the cache.
func (c *VCache) Evict() {
	c.mu.Lock()

	var evictedItems []interface{}

	for key, val := range c.store {
		if time.Now().Unix() > val.expirationTime {
			evictedItems = append(evictedItems, key)
		}
	}

	for _, key := range evictedItems {
		delete(c.store, key)
	}

	c.mu.Unlock()
}
