// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

type bucketPropertiesPageModel struct {
	P       bucketPageModel
	Details DetailsModel
}

func RenderBucketPropertiesPage(w io.Writer, b *bucket.Bucket) error {
	m := bucketPropertiesPageModel{
		P: newBucketPageModel(b, "properties"),
		Details: DetailsModel{
			Details: []DetailModel{
				{Title: "Name", Value: b.Name},
				{Title: "Objects", Value: formatInt(int(b.Objects))},
				{Title: "Size", Value: formatBytes(b.Size)},
				{Title: "Created at", Value: formatDateTime(b.CreatedAt)},
			},
		},
	}
	return renderTemplate(w, "BucketPropertiesPage", m)
}
