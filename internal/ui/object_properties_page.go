// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

type objectPropertiesPageModel struct {
	P       *bucketPageModel
	NavTabs *NavTabsModel
	Details DetailsModel
}

func RenderObjectPropertiesPage(w io.Writer, b *bucket.Bucket, o *object.Object) error {
	p := newBucketPageModel(b)
	links := NewBucketLinks(b.Name)
	p.Breadcrumbs.Last().Link = links.Objects
	addPathCrumbs(p.Breadcrumbs, links, o.Key)
	p.Breadcrumbs.Last().Link = ""
	m := objectPropertiesPageModel{
		P:       p,
		NavTabs: newObjectNavTabs(links.Folder(object.PathPrefix(o.Key, "/")), "properties"),
		Details: DetailsModel{
			Details: []DetailModel{
				{Title: "Key", Value: o.Key},
				{Title: "Size", Value: formatBytes(o.Size)},
				{Title: "Created at", Value: formatDateTime(o.CreatedAt)},
			},
		},
	}
	return renderTemplate(w, "BucketPropertiesPage", m)
}
