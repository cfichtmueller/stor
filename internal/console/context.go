// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func contextGetBucket(c jug.Context) *bucket.Bucket {
	return c.MustGet("bucket").(*bucket.Bucket)
}

func contextGetApiKey(c jug.Context) *apikey.ApiKey {
	return c.MustGet("apiKey").(*apikey.ApiKey)
}
