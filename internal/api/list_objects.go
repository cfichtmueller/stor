// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/util"
)

type ListObjectsResponse struct {
	IsTruncated    bool             `json:"isTruncated"`
	Objects        []ObjectResponse `json:"objects"`
	Name           string           `json:"name"`
	MaxKeys        int              `json:"maxKeys"`
	KeyCount       int              `json:"keyCount"`
	StartAfter     *string          `json:"startAfter,omitempty"`
	CommonPrefixes []string         `json:"commonPrefixes,omitempty"`
}

func handleListObjects(c jug.Context) {
	startAfter := c.Query("start-after")
	maxKeys, err := c.DefaultIntQuery("max-keys", 1000)
	if err != nil {
		handleError(c, err)
		return
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}
	delimiter := c.Query("delimiter")
	prefix := c.Query("prefix")
	b := contextGetBucket(c)

	if delimiter != "" {
		s := &PrefixSearch{
			b:            b,
			index:        object.NewPrefixIndex(delimiter, prefix),
			startAfter:   startAfter,
			currentStart: startAfter,
			maxKeys:      maxKeys,
			objects:      make([]*object.Object, 0),
		}
		if err := s.Do(c); err != nil {
			handleError(c, err)
			return
		}
		c.RespondOk(ListObjectsResponse{
			IsTruncated:    s.truncated,
			Objects:        util.MapMany(s.objects, newObjectResponse),
			Name:           b.Name,
			MaxKeys:        maxKeys,
			KeyCount:       len(s.objects),
			StartAfter:     &startAfter,
			CommonPrefixes: s.index.CommonPrefixes,
		})
		return
	}

	contents, err := object.List(c, b.Name, startAfter, maxKeys)
	if err != nil {
		handleError(c, err)
		return
	}
	keyCount := len(contents)
	totalKeys, err := object.Count(c, b.Name, startAfter)
	if err != nil {
		handleError(c, err)
		return
	}
	var startAfterRes *string
	if startAfter != "" {
		startAfterRes = &startAfter
	}

	c.RespondOk(ListObjectsResponse{
		IsTruncated: totalKeys > keyCount,
		Objects:     util.MapMany(contents, newObjectResponse),
		Name:        b.Name,
		MaxKeys:     maxKeys,
		KeyCount:    keyCount,
		StartAfter:  startAfterRes,
	})
}

type PrefixSearch struct {
	b            *bucket.Bucket
	index        *object.PrefixIndex
	startAfter   string
	currentStart string
	maxKeys      int
	truncated    bool
	objects      []*object.Object
}

func (s *PrefixSearch) Do(ctx context.Context) error {
	contents, err := object.List(ctx, s.b.Name, s.currentStart, 1000)
	if err != nil {
		return err
	}
	if len(contents) == 0 {
		return nil
	}
	for _, o := range contents {
		if s.index.AddKey(o.Key) {
			if len(s.objects) == s.maxKeys {
				s.truncated = true
			} else {
				s.objects = append(s.objects, o)
			}
		}
		s.currentStart = o.Key
	}
	return s.Do(ctx)
}
