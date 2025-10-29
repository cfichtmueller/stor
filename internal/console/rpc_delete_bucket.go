// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"log/slog"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/uc"
)

func handleRpcDeleteBucket(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)

	slog.Info("deleting bucket", "bucket", b.Name)

	if err := uc.DeleteBucket(c, b); err != nil {
		return srv.Respond().
			HxTrigger(hxTrigger(hxTriggerModel{
				Event: "bucketDeletionFailed",
				Toast: newToast("Error", "Failed to delete bucket: %v", err),
			}))
	}

	return srv.Respond().
		HxRedirect("/u/buckets").
		HxTrigger(hxTrigger(hxTriggerModel{
			Toast: newToast("Success", "Bucket deleted"),
		}))
}
