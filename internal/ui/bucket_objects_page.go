// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/util"
)

type BucketObjectsPageData struct {
	Bucket  *bucket.Bucket
	Objects []ObjectData
}

type bucketObjectsPageModel struct {
	P       *bucketPageModel
	NavTabs NavTabsModel
	Objects []objectModel
}

func RenderBucketObjectsPage(w io.Writer, d BucketObjectsPageData) error {
	m := &bucketObjectsPageModel{
		P:       newBucketPageModel(d.Bucket),
		NavTabs: *newBucketNavTabs(d.Bucket.Name, "objects"),
		Objects: util.MapMany(d.Objects, newObjectModel),
	}
	return renderTemplate(w, "BucketObjectsPage", m)
}
