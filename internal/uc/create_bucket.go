package uc

import (
	"context"
	"errors"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func CreateBucket(ctx context.Context, name string) (*bucket.Bucket, error) {
	exists := true
	if _, err := bucket.FindOne(ctx, name); err != nil {
		if !errors.Is(err, bucket.ErrNotFound) {
			return nil, err
		}
		exists = false
	}
	if exists {
		return nil, jug.NewConflictError("bucket exists")
	}

	b, err := bucket.Create(ctx, bucket.CreateCommand{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}
