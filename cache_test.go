package cache_test

import (
	"context"
	"testing"
	"time"

	cache "github.com/microup/vcache"
)

func TestStartEvict(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 2 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)
	err := cacheInstance.Add("test_key", "12345678")

	if err != nil {
		t.Errorf("failed add key %v", err)
	}

	cacheInstance.StartEvict(ctx)

	time.Sleep(3 * time.Second)

	if _, foundKey := cacheInstance.Get("test_key"); foundKey {
		t.Error("expected cache value to be evicted")
	}
}

func TestCacheEvict(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 1 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)
	cacheInstance.StartEvict(ctx)

	searchValue := "12345678"

	err := cacheInstance.Add("test_key", "12345678")

	if err != nil {
		t.Errorf("failed add key %v", err)
	}

	value, foundKey := cacheInstance.Get("test_key")
	if !foundKey || value != searchValue {
		t.Errorf("expected value in cache but got %v, %v", value, foundKey)
	}

	time.Sleep(2 * time.Second)

	value, foundKey = cacheInstance.Get("test_key")
	if foundKey || value != nil {
		t.Errorf("expected value to be evicted but got %v, %v", value, foundKey)
	}
}

func TestAddGetData(t *testing.T) {
	t.Parallel()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 1 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)

	searchValue := "12345678"

	err := cacheInstance.Add("test_key", searchValue)

	if err != nil {
		t.Errorf("failed add key %v", err)
	}

	value, foundKey := cacheInstance.Get("test_key")

	if !foundKey || value != searchValue {
		t.Errorf("expected value to be %s but got %s", searchValue, value)
	}

	value, foundKey = cacheInstance.Get("EMPTY")
	if foundKey || value != nil {
		t.Errorf("Expected lastUsed to be set but got nil")
	}
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	cacheInstance := cache.New(time.Minute, time.Hour)

	// Adding a key-value pair to the cache
	err := cacheInstance.Add("key1", "value1")
	
	if err != nil {
		t.Errorf("failed add key %v", err)
	}

	// Check if the key-value pair was added successfully
	val, found := cacheInstance.Get("key1")
	if !found || val != "value1" {
		t.Error("Key-value pair was not added to the cache")
	}

	// Delete the key-value pair
	cacheInstance.Delete("key1")

	// Check if the key-value pair was deleted successfully
	_, found = cacheInstance.Get("key1")
	if found {
		t.Error("Key-value pair was not deleted from the cache")
	}
}

func TestDifferentTypes(t *testing.T) {
	t.Parallel()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 10 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)

	searchValue := "12345678"

	err := cacheInstance.Add(0.75513, searchValue)

	if err != nil {
		t.Errorf("failed add key %v", err)
	}

	value, foundKey := cacheInstance.Get(0.75513)

	if !foundKey || value != searchValue {
		t.Errorf("expected value to be %s but got %s", searchValue, value)
	}

}