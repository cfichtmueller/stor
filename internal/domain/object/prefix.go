// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package object

import "strings"

type PrefixIndex struct {
	prefixIndex    map[string]bool
	CommonPrefixes []string
	keys           []string
	delimiter      string
	prefix         string
}

func NewPrefixIndex(delimiter, prefix string) *PrefixIndex {
	return &PrefixIndex{
		prefixIndex:    make(map[string]bool),
		CommonPrefixes: make([]string, 0),
		keys:           make([]string, 0),
		delimiter:      delimiter,
		prefix:         prefix,
	}
}

// AddKey adds a key to the index. The method returns true if this is a key match
func (i *PrefixIndex) AddKey(key string) bool {
	if i.prefix == "" && !strings.Contains(key, i.delimiter) {
		i.keys = append(i.keys, key)
		return true
	}
	if !strings.HasPrefix(key, i.prefix) || !strings.Contains(key, i.delimiter) {
		return false
	}
	parts := strings.Split(key[len(i.prefix):], i.delimiter)
	if len(parts) == 1 {
		i.keys = append(i.keys, i.prefix+parts[0])
		return true
	}
	commonPrefix := i.prefix + parts[0] + i.delimiter
	if _, ok := i.prefixIndex[commonPrefix]; !ok {
		i.prefixIndex[commonPrefix] = true
		i.CommonPrefixes = append(i.CommonPrefixes, commonPrefix)
	}
	return false
}
