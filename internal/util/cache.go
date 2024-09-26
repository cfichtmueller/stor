// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"sync"
	"time"
)

type Cache struct {
	index      map[string]CacheEntry
	writeMutex sync.Mutex
}

type CacheEntry struct {
	value   any
	expires time.Time
}

func NewCache() *Cache {
	return &Cache{
		index: map[string]CacheEntry{},
	}
}

func (c *Cache) Get(key string) (any, bool) {
	v, ok := c.index[key]
	if !ok {
		return nil, false
	}
	if v.expires.IsZero() || v.expires.After(time.Now()) {
		return v.value, true
	}
	c.writeMutex.Lock()
	delete(c.index, key)
	c.writeMutex.Unlock()
	return nil, false
}

func (c *Cache) SetTTL(key string, value any, ttl time.Duration) {
	c.writeMutex.Lock()
	c.index[key] = CacheEntry{
		value:   value,
		expires: time.Now().Add(ttl),
	}
	c.writeMutex.Unlock()
}
