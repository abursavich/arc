
# arc
    import "github.com/abursavich/arc"

Package arc implements an adaptive replacement cache.

See
- [Adaptive Replacement Cache](https://en.wikipedia.org/wiki/Adaptive_replacement_cache)
- [ARC: A Self-tuning, Low Overhead Replacement Cache](https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf)
- [Patent US6996676 - System and method for implementing an adaptive replacement cache policy](https://www.google.com/patents/US6996676)








## type Cache
``` go
type Cache struct {
    // contains filtered or unexported fields
}
```
Cache is an adaptive replacement cache.
It is not safe for concurrent access.









### func New
``` go
func New(size int) *Cache
```
New creates a new Cache.




### func (\*Cache) Get
``` go
func (c *Cache) Get(key Key) (value interface{}, ok bool)
```
Get reads a key's value from the cache.



### func (\*Cache) Len
``` go
func (c *Cache) Len() int
```
Len returns the number of items in the cache.



### func (\*Cache) Set
``` go
func (c *Cache) Set(key Key, value interface{})
```
Set writes a key's value to the cache.



## type Key
``` go
type Key interface{}
```
A Key may be any value that is comparable.
See [Comparison operators](http://golang.org/ref/spec#Comparison_operators).

















- - -
Partially generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)