
# arc

[![License](https://img.shields.io/badge/license-mit-blue.svg?style=for-the-badge)](https://raw.githubusercontent.com/abursavich/arc/main/LICENSE)
[![GoDev Reference](https://img.shields.io/static/v1?logo=go&logoColor=white&color=00ADD8&label=dev&message=reference&style=for-the-badge)](https://pkg.go.dev/bursavich.dev/arc)
[![Go Report Card](https://goreportcard.com/badge/bursavich.dev/arc?style=for-the-badge)](https://goreportcard.com/report/bursavich.dev/arc)

    import "bursavich.dev/arc"

Package arc implements a generic adaptive replacement cache.

See:
- [Adaptive Replacement Cache](https://en.wikipedia.org/wiki/Adaptive_replacement_cache)
- [ARC: A Self-tuning, Low Overhead Replacement Cache](https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf)
- [Patent US6996676 - System and method for implementing an adaptive replacement cache policy](https://www.google.com/patents/US6996676)


## type Cache
```go
type Cache[K comparable, V any] struct {
    // contains filtered or unexported fields
}
```
Cache is an adaptive replacement cache.
It is not safe for concurrent access.

### func New
```go
func New[K comparable, V any](size int) *Cache[K, V]
```
New creates a new Cache.

### func (\*Cache[K, V]) Get
```go
func (c *Cache[K, V]) Get(key K) (value V, ok bool)
```
Get reads a key's value from the cache.

### func (\*Cache[K, V]) Len
```go
func (c *Cache[K, V]) Len() int
```
Len returns the number of live items in the cache.


### func (\*Cache[K, V]) Set
```go
func (c *Cache[K, V]) Set(key K, value V)
```
Set writes a key's value to the cache.
