// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

var (
	paramBucketName = "bucketName"
	paramObjectKey  = "objectKey"
)

func contextGetBucket(c *srv.Context) *bucket.Bucket {
	return c.MustGet("bucket").(*bucket.Bucket)
}

func contextSetBucket(c *srv.Context, b *bucket.Bucket) {
	c.Set("bucket", b)
}

func contextGetObjectKey(c *srv.Context) string {
	return c.PathValue(paramObjectKey)
}
