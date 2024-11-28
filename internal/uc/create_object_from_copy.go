// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func CreateObjectFromCopy(ctx context.Context, b *bucket.Bucket, src *object.Object, key string) (*object.Object, error) {

	o, err := object.Copy(ctx, src, key)
	if err != nil {
		return nil, err
	}

	if err := ReconcileBucket(ctx, b); err != nil {
		return nil, err
	}

	return o, nil
}
