package cache_test

import (
	"strconv"
	"testing"
	"time"

	cache "microup.ru/vcache/src"
)

func BenchmarkCacheAdd(b *testing.B) {
	c := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		c.Add(key, value)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cacheTest := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		cacheTest.Add(key, value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cacheTest.Get(key)
	}
}

func BenchmarkCacheEvict(b *testing.B) {
	cacheTest := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		cacheTest.Add(key, value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cacheTest.Evict()
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cacheTest := cache.New(1*time.Second, 2*time.Second)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		cacheTest.Add(key, value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		cacheTest.Delete(key)
	}
}
