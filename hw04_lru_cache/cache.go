package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mu sync.RWMutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Set returns true if the key was in the cache, false otherwise.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	newItem := &cacheItem{key: key, value: value}

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = newItem
		return true
	}

	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, lastItem.Value.(*cacheItem).key)
	}

	c.items[key] = c.queue.PushFront(newItem)

	return false
}

func (c *lruCache) Clear() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.capacity = 0
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
