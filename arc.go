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

import "container/list"

// Cache is an adaptive replacement cache.
// It is not safe for concurrent access.
type Cache struct {
	n, p   int                           // max size, pivot
	rl, rd *list.List                    // MRU live, MRU dead
	fl, fd *list.List                    // MFU live, MFU dead
	tbl    map[interface{}]*list.Element // lookup table
}

// A Key may be any value that is comparable.
// See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type item struct {
	key Key
	val interface{}
	hot bool
}

const (
	miss = iota
	dead
	live
)

type sentinel struct{}

// tomb marks dead items.
var tomb interface{} = sentinel{}

// New creates a new Cache.
func New(size int) *Cache {
	if size <= 0 {
		panic("arc: size must be greater than 0")
	}
	return &Cache{
		n:  size,
		rl: list.New(), rd: list.New(),
		fl: list.New(), fd: list.New(),
		tbl: make(map[interface{}]*list.Element, size<<1),
	}
}

// Get reads a key's value from the cache.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if _, it, stat := c.get(key); stat == live {
		return it.val, true
	}
	return
}

// Set writes a key's value to the cache.
func (c *Cache) Set(key Key, value interface{}) {
	el, it, stat := c.get(key)
	switch stat {
	case live:
		it.val = value
	case dead:
		it.val = value
		if it.hot { // fd
			c.p = max(0, c.p-max(c.rd.Len()/c.fd.Len(), 1))
			c.evict(true)
			c.fd.Remove(el)
		} else { // rd
			it.hot = true
			c.p = min(c.n, c.p+max(c.fd.Len()/c.rd.Len(), 1))
			c.evict(false)
			c.rd.Remove(el)
		}
		c.tbl[key] = c.fl.PushFront(it)
	default: // miss
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
		c.tbl[key] = c.rl.PushFront(&item{key, value, false})
	}
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	return c.rl.Len() + c.fl.Len()
}

func (c *Cache) get(key Key) (el *list.Element, it *item, stat byte) {
	el = c.tbl[key]
	if el == nil {
		return nil, nil, miss
	}
	it = el.Value.(*item)
	if it.val == tomb {
		return el, it, dead
	}
	if it.hot {
		c.fl.MoveToFront(el)
	} else {
		it.hot = true
		c.rl.Remove(el)
		c.tbl[key] = c.fl.PushFront(it)
	}
	return el, it, live
}

// evict clears space by moving an item from the live cache to the dead cache.
// mfu gives preferential treatment to the MFU cache when all else is equal.
func (c *Cache) evict(mfu bool) {
	var src, dst *list.List
	if rl := c.rl.Len(); rl > 0 && (rl > c.p || (mfu && rl == c.p)) {
		src, dst = c.rl, c.rd
	} else {
		src, dst = c.fl, c.fd
	}
	e := src.Back()
	src.Remove(e)
	it := e.Value.(*item)
	it.val = tomb
	c.tbl[it.key] = dst.PushFront(it)
}

// deleteLRU removes the LRU from the list and deletes it from the table.
func (c *Cache) deleteLRU(l *list.List) {
	e := l.Back()
	l.Remove(e)
	delete(c.tbl, e.Value.(*item).key)
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
