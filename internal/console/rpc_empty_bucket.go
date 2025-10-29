// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"log/slog"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/uc"
)

func handleRpcEmptyBucket(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)

	slog.Info("emptying bucket", "bucket", b.Name)

	if err := uc.EmptyBucket(c, b); err != nil {
		return srv.Respond().
			HxTrigger(hxTrigger(hxTriggerModel{
				Toast: newToast("Error", "Failed to empty bucket: %v", err),
			}))
	}

	return srv.Respond().
		HxRefresh().
		HxTrigger(hxTrigger(hxTriggerModel{
			Toast: newToast("Success", "Bucket emptied"),
		}))
}
