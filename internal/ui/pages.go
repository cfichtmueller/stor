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
	Breadcrumbs BreadcrumbsModel
	NavTabs     NavTabsModel
	Bucket      bucketModel
	Objects     []objectModel
}

func newBucketPageModel(b *bucket.Bucket, active string) bucketPageModel {
	links := newBucketLinks(b.Name)
	return bucketPageModel{
		AppSidebar: appSidebarModel(app_sidebar_active_buckets),
		PageHeader: pageHeaderModel{
			Title:     "Buckets",
			CloseLink: bucketsLink,
		},
		Breadcrumbs: BreadcrumbsModel{
			Crumbs: []BreadcrumbModel{
				{Title: "Buckets", Link: "/u/buckets"},
				{Separator: true},
				{Title: b.Name},
			},
		},
		NavTabs: NavTabsModel{
			Tabs: []NavLink{
				{
					Link:   links.Objects,
					Active: active == "objects",
					Title:  "Objects",
					Icon:   "files",
				},
				{
					Active: active == "properties",
					Link:   links.Properties,
					Title:  "Properties",
					Icon:   "sliders-horizontal",
				},
				{
					Link:   links.Settings,
					Active: active == "settings",
					Title:  "Settings",
					Icon:   "cog",
				},
			},
		},
		Bucket: newBucketModel(b),
	}
}
