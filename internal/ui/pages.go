// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/util"
)

type bucketsPageModel struct {
	Buckets    []bucketModel
	AppSidebar SidebarModel
	PageHeader pageHeaderModel
}

func RenderBucketsPage(w io.Writer, buckets []*bucket.Bucket) error {
	return renderTemplate(w, "BucketsPage", bucketsPageModel{
		AppSidebar: appSidebarModel(app_sidebar_active_buckets),
		PageHeader: pageHeaderModel{
			Title: "Buckets",
		},
		Buckets: util.MapMany(buckets, newBucketModel),
	})
}

type bucketModel struct {
	Name    string
	Objects string
	Size    string
}

func newBucketModel(b *bucket.Bucket) bucketModel {
	return bucketModel{
		Name:    b.Name,
		Objects: formatInt(int(b.Objects)),
		Size:    formatBytes(b.Size),
	}
}

type bucketPageModel struct {
	AppSidebar  SidebarModel
	PageHeader  pageHeaderModel
	PageTitle   string
	Breadcrumbs *BreadcrumbsModel
}

func newBucketPageModel(b *bucket.Bucket) *bucketPageModel {
	return &bucketPageModel{
		AppSidebar: appSidebarModel(app_sidebar_active_buckets),
		PageHeader: pageHeaderModel{
			Title:     "Buckets",
			CloseLink: bucketsLink,
		},
		Breadcrumbs: &BreadcrumbsModel{
			Crumbs: []*BreadcrumbModel{
				{Title: "Buckets", Link: bucketsLink},
				{Separator: true},
				{Title: b.Name},
			},
		},
	}
}
