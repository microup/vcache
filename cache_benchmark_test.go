package cache_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	cache "github.com/microup/vcache"
)

func BenchmarkCacheAdd(b *testing.B) {
	c := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		_ = c.Add(key, value)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cacheTest := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		_ = cacheTest.Add(key, value)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cacheTest.Get(key)
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cacheTest := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		_ = cacheTest.Add(key, value)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cacheTest.Delete(key)
	}
}

func BenchmarkCacheEvict(b *testing.B) {
	cacheTest := cache.New(60*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		_ = cacheTest.Add(key, value)
	}

	// this needs to be done so that the condition for deletion applies to all records.
	time.Sleep(2 * time.Second)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cacheTest.Evict()
	}
}

func BenchmarkCacheMixed(b *testing.B) {
	cacheTest := cache.New(2*time.Second, 1*time.Nanosecond)
	_ = cacheTest.StartEvict(context.Background())

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		_ = cacheTest.Add(key, value)
	}
}
