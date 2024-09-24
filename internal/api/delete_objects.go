// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/util"
)

type DeleteObjectsRequest struct {
	Objects []ObjectReference `json:"objects"`
}

func (r DeleteObjectsRequest) Validate() error {
	tooMany := len(r.Objects) > 1000
	tooLittle := len(r.Objects) == 0
	v := jug.NewValidator().Require(!tooMany, "max 1000 objects are allowed").Require(!tooLittle, "at least one object is required")
	if tooMany || tooLittle {
		return v.Validate()
	}
	for i, o := range r.Objects {
		v.RequireNotEmpty(o.Key, fmt.Sprintf("object[%d].key is missing", i))
	}
	return v.Validate()
}

func handleDeleteObjects(c jug.Context) {
	b := contextGetBucket(c)
	var req DeleteObjectsRequest
	if !c.MustBindJSON(&req) {
		return
	}
	objectKeys := util.MapMany(req.Objects, func(r ObjectReference) string { return r.Key })
	objects, err := object.FindMany(c, b.Name, objectKeys, false)
	if err != nil {
		handleError(c, ec.Wrap(err))
		return
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

	b.Objects -= deletedCount
	b.Size -= deletedSize
	_ = bucket.Save(c, b)

	c.RespondOk(DeleteResults{
		Results: res,
	})
}
