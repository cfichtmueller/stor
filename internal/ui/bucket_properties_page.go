// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func BucketPropertiesPage(b *bucket.Bucket) e.Node {
	links := NewBucketLinks(b.Name)
	return BucketPage(
		links,
		bucket_navtabs_active_properties,
		PathBreadcrumbs(links, b, ""),
		PageTitle(""),
		Details("",
			Detail("Name", b.Name),
			Detail("Objects", formatInt(int(b.Objects))),
			Detail("Size", formatBytes(b.Size)),
			Detail("Created at", formatDateTime(b.CreatedAt)),
		),
	)
}
