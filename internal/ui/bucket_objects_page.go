// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

type BucketObjectsPageData struct {
	Bucket  *bucket.Bucket
	Prefix  string
	Objects []ObjectData
}

func BucketObjectsPage(d BucketObjectsPageData) e.Node {
	hasObjects := len(d.Objects) > 0
	links := NewBucketLinks(d.Bucket.Name)
	return LoggedInLayout(
		appSidebar(app_sidebar_active_buckets),
		PathBreadcrumbs(links, d.Bucket, d.Prefix),
		PageTitle(""),
		e.Div(
			e.Class("flex flex-col w-full border rounded-md bg-white"),
			BucketNavTabs(links, bucket_navtabs_active_objects),
			e.Div(
				e.Class("p-2"),
				e.Iff(hasObjects, e.F(ObjectsTable, d.Objects)),
				e.Iff(!hasObjects, BucketEmptyState),
			),
		),
	)
}
