// Copyright 2015 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package arc implements a generic adaptive replacement cache.
//
// It's similar to IBM's patented version in the sense that it pivots (or adapts) between a most-recently-used
// and frequently-used cache when misses hit a ghost cache of recently evicted entries. However, this version
// permits the deletion of entries and the precise details of pivoting rate and ghost cache eviction are different.
// Given the same sequence of mutually-supported operations, the contents of the two implementations may diverge.
//
// I am not a lawyer. This is not legal advice.
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

func (c *subCache[K, V]) remove(e *list.Element[item[K, V]]) item[K, V] {
	e.List().Remove(e)
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
	c := &Cache[K, V]{
		max:   size,
		pivot: size / 2,
	}
	c.live.init(size)
	c.dead.init(size)
	return c
}

// Len returns the number of live items in the cache.
func (c *Cache[K, V]) Len() int {
	return len(c.live.tbl)
}

// Get reads the key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, found bool) {
	if e, ok := c.get(key); ok {
		return e.Value.val, true
	}
	var zero V
	return zero, false
}

// Set writes the key's value to the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	if e, ok := c.get(key); ok {
		// Live cache hit.
		e.Value.val = value
		return
	}
	if e, ok := c.dead.tbl[key]; ok {
		// Dead cache hit.
		if e.Value.hot {
			c.pivot = max(0, c.pivot-1)
		} else {
			c.pivot = min(c.max, c.pivot+1)
		}
		c.dead.remove(e)
		c.evict(e.Value.hot)
		c.live.tbl[key] = c.live.mfu.PushFront(item[K, V]{
			key: key,
			val: value,
			hot: true,
		})
		return
	}
	// Cache miss.
	c.evict(false)
	c.live.tbl[key] = c.live.mru.PushFront(item[K, V]{
		key: key,
		val: value,
		hot: false,
	})
}

// Delete deletes the key's value from the cache.
func (c *Cache[K, V]) Delete(key K) {
	if e, ok := c.live.tbl[key]; ok {
		// Live cache hit.
		c.live.remove(e)
	} else if e, ok := c.dead.tbl[key]; ok {
		// Dead cache hit.
		c.dead.remove(e)
	}
}

func (c *Cache[K, V]) get(key K) (e *list.Element[item[K, V]], ok bool) {
	e, ok = c.live.tbl[key]
	if !ok {
		// Live cache miss.
		return nil, false
	}
	// Live cache hit.
	if e.Value.hot {
		// Key is already hot.
		c.live.mfu.MoveToFront(e)
		return e, true
	}
	// Key is newly hot.
	c.live.mru.Remove(e)
	e = c.live.mfu.PushFront(item[K, V]{
		key: key,
		val: e.Value.val,
		hot: true,
	})
	c.live.tbl[key] = e
	return e, true
}

// evict clears space, if necessary, by moving an item from the live cache to the dead cache
// and/or dropping items from the dead cache. hot gives preferential treatment to the MFU cache
// when all else is equal.
func (c *Cache[K, V]) evict(hot bool) {
	if len(c.live.tbl) >= c.max {
		mruLen := c.live.mru.Len()
		mfuLen := c.live.mfu.Len()
		live, dead := &c.live.mfu, &c.dead.mfu
		if mruLen > 0 && (mruLen > c.pivot || (hot && mruLen == c.pivot) || mfuLen == 0) {
			live, dead = &c.live.mru, &c.dead.mru
		}
		it := c.live.remove(live.Back())
		c.dead.tbl[it.key] = dead.PushFront(item[K, empty]{
			key: it.key,
			hot: it.hot,
		})
	}
	if len(c.dead.tbl) > c.max {
		dead := &c.dead.mfu
		if c.live.mru.Len()+c.dead.mru.Len() >= c.max {
			dead = &c.dead.mru
		}
		c.dead.remove(dead.Back())
	}
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
