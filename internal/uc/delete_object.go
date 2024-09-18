// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func DeleteObject(ctx context.Context, b *bucket.Bucket, o *object.Object) error {
	if err := object.Delete(ctx, o); err != nil {
		return err
	}

	b.Objects -= 1
	b.Size -= o.Size

	if err := bucket.Save(ctx, b); err != nil {
		return err
	}

	return nil
}
