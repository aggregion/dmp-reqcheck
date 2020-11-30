package common

import (
	"sync"
	"time"
)

type (
	// Cache .
	Cache interface {
		Set(key string, value interface{})
		Get(key string) (interface{}, bool)
		GetOrSet(key string, setFn func() (interface{}, error)) (interface{}, bool)
		Delete(key string)
	}

	cacheValue struct {
		ttl   int64
		value interface{}
	}

	ttlCache struct {
		mu     sync.RWMutex
		maxTTL int64
		items  map[string]interface{}
	}

	namespacedCache struct {
		orig      Cache
		namespace string
	}
)

// NewNamespacedCache .
func NewNamespacedCache(c Cache, namespace string) Cache {
	return &namespacedCache{
		orig:      c,
		namespace: namespace,
	}
}

func (c *namespacedCache) Set(key string, value interface{}) {
	c.orig.Set(c.namespace+key, value)
}

func (c *namespacedCache) Get(key string) (interface{}, bool) {
	return c.orig.Get(c.namespace + key)
}

func (c *namespacedCache) GetOrSet(key string, setFn func() (interface{}, error)) (interface{}, bool) {
	return c.orig.GetOrSet(c.namespace+key, setFn)
}

func (c *namespacedCache) Delete(key string) {
	c.orig.Delete(c.namespace + key)
}

// NewTTLCache .
func NewTTLCache(maxTTL int64) Cache {
	return &ttlCache{
		maxTTL: maxTTL,
		items:  make(map[string]interface{}),
	}
}

func (c *ttlCache) gc(someVal byte, now int64) {
	if someVal > 252 {
		c.mu.Lock()
		for key := range c.items {
			if c.items[key].(*cacheValue).ttl <= now {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *ttlCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *ttlCache) Set(key string, value interface{}) {
	now := time.Now().Unix()

	c.gc(byte(now), now)

	c.mu.Lock()
	c.items[key] = &cacheValue{
		ttl:   int64(now + int64(c.maxTTL)),
		value: value,
	}
	c.mu.Unlock()
}

func (c *ttlCache) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	var cValue interface{}
	cValue, exists = c.items[key]
	c.mu.RUnlock()
	if exists {
		now := time.Now().Unix()
		if cValue.(*cacheValue).ttl <= now {
			return nil, false
		}
		value = cValue.(*cacheValue).value
	}
	return
}

func (c *ttlCache) GetOrSet(key string, setFn func() (interface{}, error)) (value interface{}, loaded bool) {
	now := time.Now().Unix()

	c.gc(byte(now), now)

	var err error

	c.mu.Lock()
	if cValue, exists := c.items[key]; exists {
		value = cValue.(*cacheValue).value
		if cValue.(*cacheValue).ttl > now {
			loaded = true
		}
	}
	if !loaded {
		value, err = setFn()
		if err != nil {
			c.mu.Unlock()
			return nil, false
		}
		c.items[key] = &cacheValue{
			ttl:   int64(now + int64(c.maxTTL)),
			value: value,
		}
	}
	c.mu.Unlock()
	return
}
