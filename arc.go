// Copyright 2015 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.
//
// WARNING: IBM was granted a patent on the ARC algorithm.

// Package arc implements a generic adaptive replacement cache.
//
// See:
// 	https://en.wikipedia.org/wiki/Adaptive_replacement_cache
// 	https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf
// 	https://www.google.com/patents/US6996676
package arc

import "bursavich.dev/arc/internal/list"

type empty struct{}

type item[K comparable, V any] struct {
	key K
	val V
	hot bool
}

type subCache[K comparable, V any] struct {
	tbl map[K]*list.Element[item[K, V]]
	mru list.List[item[K, V]]
	mfu list.List[item[K, V]]
}

func (c *subCache[K, V]) init(size int) {
	c.tbl = make(map[K]*list.Element[item[K, V]], size)
	c.mru.Init()
	c.mfu.Init()
}

func (c *subCache[K, V]) removeBack(l *list.List[item[K, V]]) item[K, V] {
	return c.remove(l, l.Back())
}

func (c *subCache[K, V]) remove(l *list.List[item[K, V]], e *list.Element[item[K, V]]) item[K, V] {
	l.Remove(e)
	delete(c.tbl, e.Value.key)
	return e.Value
}

// Cache is an adaptive replacement cache.
// It is not safe for concurrent access.
type Cache[K comparable, V any] struct {
	max   int // max live size
	pivot int // pivot
	live  subCache[K, V]
	dead  subCache[K, empty]
}

// New creates a new Cache.
func New[K comparable, V any](size int) *Cache[K, V] {
	if size <= 0 {
		panic("arc: size must be greater than 0")
	}
	c := &Cache[K, V]{max: size}
	c.live.init(size)
	c.dead.init(size)
	return c
}

// Len returns the number of live items in the cache.
func (c *Cache[K, V]) Len() int {
	return len(c.live.tbl)
}

// Get reads a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	if e, ok := c.getLive(key); ok {
		return e.Value.val, true
	}
	return
}

// Set writes a key's value to the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	if e, ok := c.getLive(key); ok {
		// Live cache hit.
		e.Value.val = value
		return
	}
	if e, ok := c.dead.tbl[key]; ok {
		// Dead cache hit.
		if e.Value.hot {
			// MFU
			c.pivot = max(0, c.pivot-max(c.dead.mru.Len()/c.dead.mfu.Len(), 1))
			c.evict(true)
			c.dead.remove(&c.dead.mfu, e)
		} else {
			// MRU
			c.pivot = min(c.max, c.pivot+max(c.dead.mfu.Len()/c.dead.mru.Len(), 1))
			c.evict(false)
			c.dead.remove(&c.dead.mru, e)
		}
		c.live.tbl[key] = c.live.mfu.PushFront(item[K, V]{
			key: key,
			val: value,
			hot: true,
		})
		return
	}
	// Cache miss.
	if mruLen := c.live.mru.Len() + c.dead.mru.Len(); mruLen == c.max {
		if c.live.mru.Len() < c.max {
			c.dead.removeBack(&c.dead.mru)
			c.evict(false)
		} else {
			c.live.removeBack(&c.live.mru)
		}
	} else if totalSize := len(c.live.tbl) + len(c.dead.tbl); mruLen < c.max && totalSize >= c.max {
		if totalSize == c.max<<1 {
			c.dead.removeBack(&c.dead.mfu)
		}
		c.evict(false)
	}
	c.live.tbl[key] = c.live.mru.PushFront(item[K, V]{
		key: key,
		val: value,
		hot: false,
	})
}

func (c *Cache[K, V]) getLive(key K) (e *list.Element[item[K, V]], ok bool) {
	e, ok = c.live.tbl[key]
	if !ok {
		return nil, false
	}
	if e.Value.hot {
		// already hot
		c.live.mfu.MoveToFront(e)
		return e, true
	}
	// newly hot
	e.Value.hot = true
	c.live.mru.Remove(e)
	c.live.tbl[key] = c.live.mfu.PushFront(e.Value)
	return e, true
}

// evict clears space by moving an item from the live cache to the dead cache.
// mfu gives preferential treatment to the MFU cache when all else is equal.
func (c *Cache[K, V]) evict(mfu bool) {
	live, dead := &c.live.mfu, &c.dead.mfu
	if n := c.live.mru.Len(); n > 0 && (n > c.pivot || (mfu && n == c.pivot)) {
		live, dead = &c.live.mru, &c.dead.mru
	}
	it := c.live.removeBack(live)
	c.dead.tbl[it.key] = dead.PushFront(item[K, empty]{
		key: it.key,
		hot: it.hot,
	})
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
