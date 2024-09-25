// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func ReconcileBucket(ctx context.Context, b *bucket.Bucket) error {
	s, err := object.StatsForBucket(ctx, b.Name)
	if err != nil {
		return err
	}

	b.Objects = s.ObjectCount
	b.Size = s.TotalSize

	if err := bucket.Save(ctx, b); err != nil {
		return err
	}

	return nil
}
