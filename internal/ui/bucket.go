// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func RenderBucketsEmptyState(w io.Writer) error {
	return renderTemplate(w, "BucketsEmptyState", nil)
}

func RenderBucketsTable(w io.Writer, buckets []*bucket.Bucket) error {
	return renderTemplate(w, "BucketsTable", map[string]any{
		"Buckets": buckets,
	})
}

func RenderCreateBucketDialog(w io.Writer) error {
	return renderTemplate(w, "CreateBucketDialog", nil)
}
