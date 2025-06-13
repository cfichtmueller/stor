// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/cfichtmueller/srv"
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

func handleBucketPost(c *srv.Context) *srv.Response {
	q := c.Request().URL.Query()
	if q.Has("delete") {
		return handleDeleteObjects(c)
	}
	return srv.Respond().MethodNotAllowed()
}

func handleCreateBucket(c *srv.Context) *srv.Response {
	name := c.PathValue("bucketName")
	if err := bucket.ValidateName(name); err != nil {
		return responseFromError(err)
	}

	b, err := uc.CreateBucket(c, name)
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().Created(newBucketResponse(b))
}

func handleDeleteBucket(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)

	count, err := object.Count(c, b.Name, "")
	if err != nil {
		return responseFromError(err)
	}
	if count > 0 {
		return responseFromError(ec.BucketNotEmpty)
	}

	bucket.Delete(c, b.Name)

	return srv.Respond().NoContent()
}
