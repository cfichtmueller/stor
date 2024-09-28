// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

type ObjectPrefixSearchResult struct {
	IsTruncated    bool
	CommonPrefixes []string
	Objects        []*object.Object
}

func ObjectPrefixSearch(ctx context.Context, b *bucket.Bucket, delimiter, prefix, startAfter string, maxKeys int) (*ObjectPrefixSearchResult, error) {
	s := &objectPrefixSearch{
		b:            b,
		index:        object.NewPrefixIndex(delimiter, prefix),
		startAfter:   startAfter,
		currentStart: startAfter,
		maxKeys:      maxKeys,
		objects:      make([]*object.Object, 0),
	}
	if err := s.Do(ctx); err != nil {
		return nil, err
	}
	return &ObjectPrefixSearchResult{
		IsTruncated:    s.truncated,
		CommonPrefixes: s.index.CommonPrefixes,
		Objects:        s.objects,
	}, nil
}

type objectPrefixSearch struct {
	b            *bucket.Bucket
	index        *object.PrefixIndex
	startAfter   string
	currentStart string
	maxKeys      int
	truncated    bool
	objects      []*object.Object
}

func (s *objectPrefixSearch) Do(ctx context.Context) error {
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
