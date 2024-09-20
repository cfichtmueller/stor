// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import "time"

type Cache struct {
	index map[string]CacheEntry
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
	delete(c.index, key)
	return nil, false
}

func (c *Cache) SetTTL(key string, value any, ttl time.Duration) {
	c.index[key] = CacheEntry{
		value:   value,
		expires: time.Now().Add(ttl),
	}
}
