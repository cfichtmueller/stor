package api

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/util"
)

type ListBucketResponse struct {
	Buckets     []BucketResponse `json:"buckets"`
	IsTruncated bool             `json:"isTruncated"`
}

func handleListBuckets(c jug.Context) {
	startAfter := c.Query("start-after")
	maxBuckets, err := c.DefaultIntQuery("max-buckets", 1000)
	if err != nil {
		handleError(c, ec.InvalidArgument)
		return
	}

	buckets, err := bucket.List(c, startAfter, maxBuckets)
	if err != nil {
		handleError(c, err)
		return
	}
	count, err := bucket.Count(c, startAfter)
	if err != nil {
		handleError(c, err)
		return
	}

	c.RespondOk(ListBucketResponse{
		Buckets:     util.MapMany(buckets, newBucketResponse),
		IsTruncated: count > len(buckets),
	})
}
