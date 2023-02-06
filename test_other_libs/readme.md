## Results benchmark on cpu: AMD Ryzen 5 5600X 6-Core Processor

```
go test -bench=. -benchmem -benchtime=5s
VCacheAdd-12             8966773               761.1 ns/op           213 B/op          7 allocs/op
VCacheGet-12            39050804               184.0 ns/op             7 B/op          0 allocs/op
VCacheDelete-12         36885169               192.3 ns/op             7 B/op          0 allocs/op
VCacheEvict-12          23756037               248.5 ns/op             0 B/op          0 allocs/op
VCacheMixed-12           9602469               791.1 ns/op           154 B/op          7 allocs/op
```

The results of comparison with another librarys:

- [go-cache](https://github.com/patrickmn/go-cache) are also presented:

```
go test -bench=. -benchmem -benchtime=5s
GoCacheAdd-12            11722840               459.1 ns/op           203 B/op          5 allocs/op
GoCacheGet-12            40415334               185.9 ns/op             7 B/op          0 allocs/op
GoCacheDelete-12         35583009               218.9 ns/op            23 B/op          1 allocs/op
```

- [go-generics-cache](https://github.com/Code-Hex/go-generics-cache) are also presented:

```
GoCGenericsCacheAdd-12      7212994             745.5 ns/op           304 B/op          7 allocs/op
GoCGenericsCacheGet-12      31680290            201.4 ns/op            23 B/op          1 allocs/op
GoCGenericsCacheDelete-12   36424033            231.1 ns/op            23 B/op          1 allocs/op
```
