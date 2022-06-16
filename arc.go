// Copyright 2015 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.
//
// WARNING: IBM was granted a patent on the ARC algorithm.

// Package arc implements an adaptive replacement cache.
//
// See
// 	https://en.wikipedia.org/wiki/Adaptive_replacement_cache
// 	https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf
// 	https://www.google.com/patents/US6996676
package arc

import "bursavich.dev/arc/internal/list"

const (
	live = 1 << iota
	hot
)

type item[K comparable, V any] struct {
	key  K
	val  V
	mask int
}

func (i *item[K, V]) has(v int) bool { return i.mask&v == v }
func (i *item[K, V]) set(v int)      { i.mask |= v }
func (i *item[K, V]) unset(v int)    { i.mask &= ^v }

// Cache is an adaptive replacement cache.
// It is not safe for concurrent access.
type Cache[K comparable, V any] struct {
	n, p   int                             // max size, pivot
	rl, rd *list.List[item[K, V]]          // MRU live, MRU dead
	fl, fd *list.List[item[K, V]]          // MFU live, MFU dead
	tbl    map[K]*list.Element[item[K, V]] // lookup table
}

// New creates a new Cache.
func New[K comparable, V any](size int) *Cache[K, V] {
	if size <= 0 {
		panic("arc: size must be greater than 0")
	}
	return &Cache[K, V]{
		n:  size,
		rl: list.New[item[K, V]](), rd: list.New[item[K, V]](),
		fl: list.New[item[K, V]](), fd: list.New[item[K, V]](),
		tbl: make(map[K]*list.Element[item[K, V]], size<<1),
	}
}

// Get reads a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	if el, ok := c.get(key); ok && el.Value.has(live) {
		return el.Value.val, true
	}
	return
}

// Set writes a key's value to the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	if el, ok := c.get(key); !ok {
		// miss
		if l1 := c.rl.Len() + c.rd.Len(); l1 == c.n {
			if c.rl.Len() < c.n {
				c.deleteLRU(c.rd)
				c.evict(false)
			} else {
				c.deleteLRU(c.rl)
			}
		} else if l1 < c.n && len(c.tbl) >= c.n {
			if len(c.tbl) == c.n<<1 {
				c.deleteLRU(c.fd)
			}
			c.evict(false)
		}
		c.tbl[key] = c.rl.PushFront(item[K, V]{key, value, live})
	} else if el.Value.has(live) {
		// live
		el.Value.val = value
	} else {
		// dead
		el.Value.val = value
		if el.Value.has(hot) { // fd
			c.p = max(0, c.p-max(c.rd.Len()/c.fd.Len(), 1))
			c.evict(true)
			c.fd.Remove(el)
		} else { // rd
			el.Value.set(hot)
			c.p = min(c.n, c.p+max(c.fd.Len()/c.rd.Len(), 1))
			c.evict(false)
			c.rd.Remove(el)
		}
		el.Value.set(live)
		c.tbl[key] = c.fl.PushFront(el.Value)
	}
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	return c.rl.Len() + c.fl.Len()
}

func (c *Cache[K, V]) get(key K) (el *list.Element[item[K, V]], ok bool) {
	el = c.tbl[key]
	if el == nil {
		return nil, false
	}
	if !el.Value.has(live) { // dead
		return el, true
	}
	if el.Value.has(hot) { // hot
		c.fl.MoveToFront(el)
		return el, true
	}
	// live
	el.Value.set(hot)
	c.rl.Remove(el)
	c.tbl[key] = c.fl.PushFront(el.Value)
	return el, true
}

// evict clears space by moving an item from the live cache to the dead cache.
// mfu gives preferential treatment to the MFU cache when all else is equal.
func (c *Cache[K, V]) evict(mfu bool) {
	var src, dst *list.List[item[K, V]]
	if rl := c.rl.Len(); rl > 0 && (rl > c.p || (mfu && rl == c.p)) {
		src, dst = c.rl, c.rd
	} else {
		src, dst = c.fl, c.fd
	}
	el := src.Back()
	src.Remove(el)
	el.Value.unset(live)
	c.tbl[el.Value.key] = dst.PushFront(el.Value)
}

// deleteLRU removes the LRU from the list and deletes it from the table.
func (c *Cache[K, V]) deleteLRU(l *list.List[item[K, V]]) {
	e := l.Back()
	l.Remove(e)
	delete(c.tbl, e.Value.key)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
