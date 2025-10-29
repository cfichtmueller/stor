// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func DeleteBucketDialog(b *bucket.Bucket) e.Node {
	return ConfirmationDialog(&ConfirmationDialogData{
		ID:           "delete-bucket-dialog",
		Title:        "Delete Bucket",
		Description:  "Are you sure you want to delete this bucket? This action cannot be undone.",
		ConfirmTitle: "Delete Bucket",
		CancelTitle:  "Cancel",
		Destructive:  true,
		HxDelete:     "/r/bucket?bucket=" + b.Name,
	})
}
