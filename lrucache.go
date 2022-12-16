package lrucache

import "sync"

type LRUCache[K comparable] struct {
	mapping  map[K]*LRUCacheEntry[K]
	tail     *LRUCacheEntry[K]
	head     *LRUCacheEntry[K]
	size     int
	capacity int
	lock     sync.Mutex
}

type LRUCacheEntry[K comparable] struct {
	key      K
	value    any
	next     *LRUCacheEntry[K]
	previous *LRUCacheEntry[K]
}

func NewLRUCache[K comparable](capacity int) LRUCache[K] {
	return LRUCache[K]{
		mapping:  make(map[K]*LRUCacheEntry[K], capacity),
		tail:     nil,
		head:     nil,
		size:     0,
		capacity: capacity,
	}
}

func newLRUCacheEntry[K comparable](key K, value any) *LRUCacheEntry[K] {
	return &LRUCacheEntry[K]{
		key:      key,
		value:    value,
		next:     nil,
		previous: nil,
	}
}

func (c *LRUCache[K]) removeEntry(entry *LRUCacheEntry[K]) {
	if entry == c.head {
		c.head = c.head.previous
		c.head.next = nil
	} else if entry == c.tail {
		c.tail = c.tail.next
		c.tail.previous = nil
	} else {
		entry.previous.next = entry.next
		entry.next.previous = entry.previous
	}
}

func (c *LRUCache[K]) addEntry(entry *LRUCacheEntry[K]) {
	if c.head != nil {
		c.head.next = entry
		entry.previous = c.head
		c.head = entry
	} else {
		c.head = entry
		c.tail = entry
	}
	entry.next = nil
	c.mapping[entry.key] = entry
}

func (c *LRUCache[K]) moveToFront(entry *LRUCacheEntry[K]) {
	if entry == c.head {
		return
	}

	c.removeEntry(entry)
	c.addEntry(entry)
}

func (c *LRUCache[K]) Get(key K) any {
	c.lock.Lock()
	defer c.lock.Unlock()
	if entry, ok := c.mapping[key]; ok {
		c.moveToFront(entry)
		return entry.value
	}
	return nil
}

func (c *LRUCache[K]) Set(key K, value any) {
	c.lock.Lock()
	if entry, ok := c.mapping[key]; ok {
		c.moveToFront(entry)
		return
	}

	entry := newLRUCacheEntry(key, value)
	if c.size >= c.capacity {
		delete(c.mapping, c.tail.key)
		nextEntry := c.tail.next
		nextEntry.previous = nil
		c.tail = nextEntry
	} else {
		c.size += 1
	}

	c.addEntry(entry)
	c.lock.Unlock()
}
