// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func contextGetBucket(c jug.Context) *bucket.Bucket {
	return c.MustGet("bucket").(*bucket.Bucket)
}

func contextGetObjectKey(c jug.Context) string {
	return c.Param("objectKey")[1:]
}

func contextGetObject(c jug.Context) *object.Object {
	return c.MustGet("object").(*object.Object)
}
