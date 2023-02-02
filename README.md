[![test-and-linter](https://github.com/microup/vcache/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/microup/vcache/actions/workflows/main.yml)

# What is VCache?

This is a Go package named "vcache" which implements a simple in-memory cache. The cache stores data as key-value pairs, where the keys are interface types and the values can be of any type. The cache is designed to be concurrent-safe with the use of a sync.Mutex.

* This library "vcache" is a pure implementation and does not rely on any external dependencies. It is a self-contained implementation of an in-memory cache.

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

## Results benchmark

```
go test -bench=. -benchmem -benchtime=5s

cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkCacheAdd-12             7467901               873.1 ns/op           259 B/op          7 allocs/op
BenchmarkCacheGet-12            28955613               334.7 ns/op            39 B/op          2 allocs/op
BenchmarkCacheDelete-12         36772456               204.5 ns/op             7 B/op          0 allocs/op
BenchmarkCacheEvict-12          22011206               259.6 ns/op             0 B/op          0 allocs/op
```
