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

type BucketFolderPageData struct {
	Bucket  *bucket.Bucket
	Prefix  string
	Objects []ObjectData
}

type bucketFolderPageModel struct {
	P       *bucketPageModel
	NavTabs *NavTabsModel
	Objects []objectModel
}

func RenderBucketFolderPage(w io.Writer, d BucketFolderPageData) error {
	p := newBucketPageModel(d.Bucket)
	links := NewBucketLinks(d.Bucket.Name)
	p.Breadcrumbs.Last().Link = links.Objects
	prefix := ""
	for _, f := range object.SplitPath(d.Prefix, "/") {
		prefix = prefix + f + "/"
		p.Breadcrumbs.Add(&BreadcrumbModel{Separator: true})
		p.Breadcrumbs.Add(&BreadcrumbModel{
			Title: f,
			Link:  links.Folder(prefix),
		})
	}
	p.Breadcrumbs.Last().Link = ""
	m := &bucketFolderPageModel{
		P:       p,
		NavTabs: newBuckeFoldertNavTabs(d.Bucket.Name),
		Objects: util.MapMany(d.Objects, newObjectModel),
	}
	return renderTemplate(w, "BucketObjectsPage", m)
}
