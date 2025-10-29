// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func EmptyBucketDialog(b *bucket.Bucket) e.Node {
	return ConfirmationDialog(&ConfirmationDialogData{
		ID:           "empty-bucket-dialog",
		Title:        "Empty Bucket",
		Description:  "Are you sure you want to empty this bucket? This action cannot be undone.",
		ConfirmTitle: "Empty Bucket",
		CancelTitle:  "Cancel",
		Destructive:  true,
		HxPost:       "/r/empty-bucket?bucket=" + b.Name,
	})
}
