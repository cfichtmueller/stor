// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func EmptyBucket(ctx context.Context, b *bucket.Bucket) error {
	var deleteObjects func() error

	deleteObjects = func() error {
		objects, err := object.List(ctx, b.Name, "", 1000)
		if err != nil {
			return err
		}

		if len(objects) == 0 {
			return nil
		}

		for _, o := range objects {
			if err := object.Delete(ctx, o); err != nil {
				return err
			}
		}

		return deleteObjects()
	}

	if err := deleteObjects(); err != nil {
		return err
	}

	if err := ReconcileBucket(ctx, b); err != nil {
		return err
	}

	return nil
}
