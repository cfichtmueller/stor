// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log/slog"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/util"
)

type DeleteObjectsRequest struct {
	Objects []ObjectReference `json:"objects"`
}

func (r DeleteObjectsRequest) Validate() error {
	v := srv.RequireMinLengthSlice("objects", 1, r.Objects, nil)
	v = srv.RequireMaxLengthSlice("objects", 1000, r.Objects, v)
	if v != nil {
		return v
	}

	for i, o := range r.Objects {
		v = srv.RequireNotEmptyIndexed("objects[%d].key", i, o.Key, v)
	}
	return srv.Validate(v)
}

func handleDeleteObjects(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	var req DeleteObjectsRequest
	if r := c.BindJSON(&req); r != nil {
		return r
	}
	objectKeys := util.MapMany(req.Objects, func(r ObjectReference) string { return r.Key })
	objects, err := object.FindMany(c, b.Name, objectKeys, false)
	if err != nil {
		return responseFromError(ec.Wrap(err))
	}
	var deletedCount int64 = 0
	var deletedSize int64 = 0
	index := make(map[string]*DeleteResult)
	for _, o := range objects {
		res := &DeleteResult{
			Key: o.Key,
		}
		if err := object.Delete(c, o); err != nil {
			res.Error = ec.Wrap(err)
		} else {
			res.Deleted = true
			deletedCount += 1
			deletedSize += o.Size
		}
		index[o.Key] = res
	}
	for _, o := range req.Objects {
		if _, ok := index[o.Key]; !ok {
			index[o.Key] = &DeleteResult{
				Key:     o.Key,
				Deleted: false,
				Error:   ec.NoSuchKey,
			}
		}
	}

	res := make([]DeleteResult, 0, len(index))
	for _, r := range index {
		res = append(res, *r)
	}

	if err := uc.ReconcileBucket(c, b); err != nil {
		slog.Error("unable to reconcile bucket", "error", err)
	}

	return srv.Respond().Json(DeleteResults{
		Results: res,
	})
}
