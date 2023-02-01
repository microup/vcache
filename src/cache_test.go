package cache_test

import (
	"context"
	"testing"
	"time"

	"microup.ru/vcache/src"
)

func TestStartEvict(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 2 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)
	cacheInstance.Add("test_key", "12345678")
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

	cacheInstance.Add("test_key", "12345678")

	value, foundKey := cacheInstance.Get("test_key")
	if !foundKey || value != searchValue {
		t.Errorf("expected value in cache but got %v, %v", value, foundKey)
	}

	time.Sleep(2 * time.Second)

	value, foundKey = cacheInstance.Get("test_key")
	if foundKey || value != "" {
		t.Errorf("expected value to be evicted but got %v, %v", value, foundKey)
	}
}

func TestAddGetData(t *testing.T) {
	t.Parallel()

	timeCheckNewTicker := 1 * time.Second
	timeRecordEvict := 1 * time.Second

	cacheInstance := cache.New(timeCheckNewTicker, timeRecordEvict)

	searchValue := "12345678"

	cacheInstance.Add("test_key", searchValue)
	value, foundKey := cacheInstance.Get("test_key")

	if !foundKey || value != searchValue {
		t.Errorf("expected value to be %s but got %s", searchValue, value)
	}

	value, foundKey = cacheInstance.Get("EMPTY")
	if foundKey || value != "" {
		t.Errorf("Expected lastUsed to be set but got nil")
	}
}
