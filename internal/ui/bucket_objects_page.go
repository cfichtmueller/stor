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

func RenderBucketObjectsPage(w io.Writer, b *bucket.Bucket, objects []*object.Object) error {
	m := newBucketPageModel(b, "objects")
	m.Objects = util.MapMany(objects, newObjectModel)
	return renderTemplate(w, "BucketObjectsPage", m)
}
