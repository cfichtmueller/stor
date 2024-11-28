// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func CreateObjectFromData(ctx context.Context, b *bucket.Bucket, cmd object.CreateCommand) (*object.Object, error) {
	exists, err := object.Exists(ctx, b.Name, cmd.Key)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, jug.NewConflictError("object exists")
	}

	o, err := object.Create(ctx, b.Name, cmd)
	if err != nil {
		return nil, err
	}

	if err := ReconcileBucket(ctx, b); err != nil {
		return nil, err
	}

	return o, nil
}
