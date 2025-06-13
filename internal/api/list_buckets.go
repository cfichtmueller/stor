package api

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/util"
)

type ListBucketResponse struct {
	Buckets     []BucketResponse `json:"buckets"`
	IsTruncated bool             `json:"isTruncated"`
}

func handleListBuckets(c *srv.Context) *srv.Response {
	startAfter := c.Query("start-after")
	maxBuckets, r := c.IntQueryOrDefault("max-buckets", 1000)
	if r != nil {
		return r
	}

	buckets, err := bucket.List(c, startAfter, maxBuckets)
	if err != nil {
		return responseFromError(err)
	}
	count, err := bucket.Count(c, startAfter)
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().Json(ListBucketResponse{
		Buckets:     util.MapMany(buckets, newBucketResponse),
		IsTruncated: count > len(buckets),
	})
}
