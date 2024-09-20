package api

import (
	"fmt"
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/util"
)

type BucketResponse struct {
	Name      string    `json:"name"`
	Objects   uint64    `json:"objects"`
	Size      uint64    `json:"size"`
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

type BulkDeleteRequest struct {
	Objects []ObjectReference `json:"objects`
}

func (r BulkDeleteRequest) Validate() error {
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

type ObjectReference struct {
	Key string `json:"key"`
}

type DeleteResults struct {
	Results []DeleteResult `json:"results"`
}

type DeleteResult struct {
	Key     string `json:"key"`
	Deleted bool   `json:"deleted"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string
	Message string
}

func handleBucketPost(c jug.Context) {
	q := c.Request().URL.Query()
	if q.Has("delete") {
		handleBulkDelete(c)
		return
	}
	c.Status(405)
}

func handleBulkDelete(c jug.Context) {
	b := contextGetBucket(c)
	var req BulkDeleteRequest
	if !c.MustBindJSON(&req) {
		return
	}
	objectKeys := util.MapMany(req.Objects, func(r ObjectReference) string { return r.Key })
	objects, err := object.FindMany(c, b.Name, objectKeys)
	if err != nil {
		c.HandleError(err)
		return
	}
	deletedCount := 0
	var deletedSize uint64 = 0
	index := make(map[string]*DeleteResult)
	for _, o := range objects {
		res := &DeleteResult{
			Key: o.Key,
		}
		if err := object.Delete(c, o); err != nil {
			//TODO: all errors should have a code
			res.Error = &Error{
				Code:    "",
				Message: err.Error(),
			}
			res.Deleted = false
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
				Error: &Error{
					Code:    "ObjectNotFound",
					Message: "Object not found",
				},
			}
		}
	}

	res := make([]DeleteResult, 0, len(index))
	for _, r := range index {
		res = append(res, *r)
	}

	b.Objects -= uint64(deletedCount)
	b.Size -= deletedSize
	_ = bucket.Save(c, b)

	c.RespondOk(DeleteResults{
		Results: res,
	})
}

func handleCreateBucket(c jug.Context) {
	name := c.Param("bucketName")

	b, err := uc.CreateBucket(c, name)
	if err != nil {
		c.HandleError(err)
		return
	}

	c.RespondCreated(newBucketResponse(b))
}
