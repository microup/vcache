[![test-and-linter](https://github.com/microup/vcache/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/microup/vcache/actions/workflows/main.yml)

# What is VCache?

This is a Go package named "vcache" which implements a simple in-memory cache. The cache stores data as key-value pairs, where the keys are interface types and the values can be of any type. The cache is designed to be concurrent-safe with the use of a sync.Mutex.

* This library "vcache" is a pure implementation and does not rely on any external dependencies. It is a self-contained implementation of an in-memory cache.

This project differs from the more well-known [go-cache](https://github.com/patrickmn/go-cache) in that it uses an intraface{} type key instead of a string key in its map structure, making cache management more flexible.

## Where can this be used?

This library can be applied in software systems where caching is needed, such as web applications, databases, or other systems that require fast access to frequently used data. The library provides an in-memory cache that can be used to store frequently used data, such as API responses, database results, or other frequently accessed data. By using this library, the system can avoid unnecessary data processing and improve performance by quickly retrieving the data from the cache. The cache can also be automatically cleaned up based on the time specified in the durationTimeEvict variable, freeing up memory and ensuring that the cache remains relevant.

## It has the following main functions:

- New() creates a new instance of the cache. The parameters are timeCheckNewTicker and timeRecordEvict, which represent the frequency of cache eviction check and the time an entry can stay in the cache before being evicted, respectively.
- StartEvict() starts the cache eviction process in a separate go routine.
- Add() adds a key-value pair to the cache.
- Get() retrieves the value for a given key in the cache.
- Evict() is responsible for removing stale cache entries that have been in the cache for a duration of time greater than durationTimeEvict.

## Automatic Cache Eviction in "vcache" Library

The cache entries are automatically evicted every durationCheckTicker time intervals. The eviction process removes any cache entries that have not been accessed within the specified durationTimeEvict time period.

## Here is an example of how to use the cache library

```go
package main

import (
    "context"
    "fmt"
    "time"

    "cache"
)

func main() {
    // Create a context for the cache eviction
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Create a cache with a check ticker of 1 second and a record eviction of 5 seconds
    c := cache.New(time.Second, 5*time.Second)

    // Start eviction routine
    c.StartEvict(ctx)

    // Add key-value pairs to the cache
    err := c.Add("key1", "value1")
    if err != nil {
        fmt.Println(err)
    }

    err = c.Add("key2", 2)
    if err != nil {
        fmt.Println(err)
    }

    // Get values from the cache
    val, found := c.Get("key1")
    if found {
        fmt.Printf("key1: %v\n", val)
    } else {
        fmt.Println("key1 not found")
    }

    val, found = c.Get("key2")
    if found {
        fmt.Printf("key2: %v\n", val)
    } else {
        fmt.Println("key2 not found")
    }

    // Wait for 6 seconds for key1 to be evicted
    time.Sleep(6 * time.Second)

    // Try to get key1 after eviction
    val, found = c.Get("key1")
    if found {
        fmt.Printf("key1: %v\n", val)
    } else {
        fmt.Println("key1 not found")
    }
}
```
output:

```bash
key1: value1
key2: 2
key1 not found
```

In this example, the cache is created with a check ticker of 1 second and a record eviction of 5 seconds. The cache eviction routine is started using the StartEvict method and passing in the context created earlier. Key-value pairs are added to the cache using the Add method, and values are retrieved using the Get method. After waiting for 6 seconds, the Get method is used again to retrieve the value for the key "key1", but it is no longer found because it has been evicted from the cache.

## Results benchmark

```
go test -bench=. -benchmem -benchtime=5s

cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkCacheAdd-12             9367170               759.6 ns/op           232 B/op          7 allocs/op
BenchmarkCacheGet-12            39035611               194.4 ns/op             7 B/op          0 allocs/op
BenchmarkCacheEvict-12          22821903               258.2 ns/op             0 B/op          0 allocs/op
BenchmarkCacheDelete-12         37573690               194.3 ns/op             7 B/op          0 allocs/op
```
