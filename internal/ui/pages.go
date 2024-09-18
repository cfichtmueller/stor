// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
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
	Bucket     bucketModel
	AppSidebar SidebarModel
	PageHeader pageHeaderModel
	NavTabs    NavTabsModel
	Objects    []objectModel
}

func newBucketPageModel(b *bucket.Bucket, active string) bucketPageModel {
	links := newBucketLinks(b.Name)
	return bucketPageModel{
		AppSidebar: appSidebarModel(app_sidebar_active_buckets),
		PageHeader: pageHeaderModel{
			Title:     "Buckets",
			CloseLink: bucketsLink,
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

func RenderBucketPage(w io.Writer, b *bucket.Bucket) error {
	return renderTemplate(w, "BucketPage", newBucketPageModel(b, "files"))
}

func RenderBucketObjectsPage(w io.Writer, b *bucket.Bucket, objects []*object.Object) error {
	m := newBucketPageModel(b, "objects")
	m.Objects = util.MapMany(objects, newObjectModel)
	return renderTemplate(w, "BucketObjectsPage", m)
}

func RenderBucketSettingsPage(w io.Writer, b *bucket.Bucket) error {
	return renderTemplate(w, "BucketSettingsPage", newBucketPageModel(b, "settings"))
}
