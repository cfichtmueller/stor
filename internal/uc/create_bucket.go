package uc

import (
	"context"
	"errors"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/ec"
)

func CreateBucket(ctx context.Context, name string) (*bucket.Bucket, error) {
	exists := true
	if _, err := bucket.FindOne(ctx, name); err != nil {
		if !errors.Is(err, ec.NoSuchBucket) {
			return nil, err
		}
		exists = false
	}
	if exists {
		return nil, ec.BucketAlreadyExists
	}

	b, err := bucket.Create(ctx, bucket.CreateCommand{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}
