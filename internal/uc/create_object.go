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

func CreateObject(ctx context.Context, b *bucket.Bucket, cmd object.CreateCommand) error {
	exists, err := object.Exists(ctx, b.Name, cmd.Key)
	if err != nil {
		return err
	}

	if exists {
		return jug.NewConflictError("object exists")
	}

	if err := object.Create(ctx, b.Name, cmd); err != nil {
		return err
	}

	b.AddObject(uint64(len(cmd.Data)))

	if err := bucket.Save(ctx, b); err != nil {
		return err
	}
	return nil
}