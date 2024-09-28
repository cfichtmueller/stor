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

func newBucketNavTabs(bucketName, active string) *NavTabsModel {
	links := NewBucketLinks(bucketName)
	return &NavTabsModel{
		Tabs: []*NavLink{
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
	}
}

func newBuckeFoldertNavTabs(bucketName string) *NavTabsModel {
	return &NavTabsModel{
		Tabs: []*NavLink{
			{
				Active: true,
				Title:  "Objects",
				Icon:   "files",
			},
		},
	}
}
