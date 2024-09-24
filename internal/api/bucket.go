// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
)

type BucketResponse struct {
	Name      string    `json:"name"`
	Objects   int64     `json:"objects"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

func newBucketResponse(b *bucket.Bucket) BucketResponse {
	return BucketResponse{
		Name:      b.Name,
		Objects:   b.Objects,
		Size:      b.Size,
		CreatedAt: b.CreatedAt,
	}
}

type ObjectReference struct {
	Key string `json:"key"`
}

type DeleteResults struct {
	Results []DeleteResult `json:"results"`
}

type DeleteResult struct {
	Key     string    `json:"key"`
	Deleted bool      `json:"deleted"`
	Error   *ec.Error `json:"error,omitempty"`
}

func handleBucketPost(c jug.Context) {
	q := c.Request().URL.Query()
	if q.Has("delete") {
		handleDeleteObjects(c)
		return
	}
	c.Status(405)
}

func handleCreateBucket(c jug.Context) {
	name := c.Param("bucketName")

	b, err := uc.CreateBucket(c, name)
	if err != nil {
		handleError(c, err)
		return
	}

	c.RespondCreated(newBucketResponse(b))
}

func handleDeleteBucket(c jug.Context) {
	b := contextGetBucket(c)

	count, err := object.Count(c, b.Name, "")
	if err != nil {
		handleError(c, err)
		return
	}
	if count > 0 {
		handleError(c, ec.BucketNotEmpty)
		return
	}

	bucket.Delete(c, b.Name)

	c.RespondNoContent()
}
