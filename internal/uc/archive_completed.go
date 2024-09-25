// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/archive"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func onArchiveCompleted(ctx context.Context, e any) error {
	d := e.(archive.CompletedEvent)

	b, err := bucket.FindOne(ctx, d.Bucket)
	if err != nil {
		return err
	}

	if err := ReconcileBucket(ctx, b); err != nil {
		return err
	}

	return nil
}
