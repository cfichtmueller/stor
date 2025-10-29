// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
)

func DeleteBucket(ctx context.Context, b *bucket.Bucket) error {
	count, err := object.Count(ctx, b.Name, "")
	if err != nil {
		return err
	}
	if count > 0 {
		return ec.BucketNotEmpty
	}

	if err := bucket.Delete(ctx, b.Name); err != nil {
		return err
	}

	return nil
}
