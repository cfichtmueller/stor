// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func contextSetBucket(c *srv.Context, b *bucket.Bucket) {
	c.Set("bucket", b)
}

func contextGetBucket(c *srv.Context) *bucket.Bucket {
	return c.MustGet("bucket").(*bucket.Bucket)
}

func contextGetApiKey(c *srv.Context) *apikey.ApiKey {
	return c.MustGet("apiKey").(*apikey.ApiKey)
}

func contextSetPrincipal(c *srv.Context, principal string) {
	c.Set("principal", principal)
}

func contextMustGetPrincipal(c *srv.Context) string {
	return c.MustGet("principal").(string)
}
