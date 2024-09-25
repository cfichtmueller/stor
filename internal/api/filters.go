// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"strings"
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/archive"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/util"
)

var (
	tokenCache = util.NewCache()
	tokenTTL   = time.Minute
)

func bucketFilter(c jug.Context) {
	name := c.Param("bucketName")
	if len(name) < 3 {
		handleError(c, ec.NoSuchBucket)
		c.Abort()
		return
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		handleError(c, err)
		c.Abort()
		return
	}
	c.Set("bucket", b)
}

func objectFilter(c jug.Context) (*object.Object, bool) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)

	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil {
		handleError(c, err)
		return nil, false
	}

	return o, true
}

func authenticatedFilter(c jug.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		handleError(c, ec.Unauthorized)
		c.Abort()
		return
	}
	token := auth[7:]
	if p, ok := tokenCache.Get(token); ok {
		c.Set("principal", p)
		c.Next()
		return
	}
	key, err := apikey.Authenticate(c, token)
	if err != nil {
		handleError(c, err)
		c.Abort()
		return
	}
	principal := "apikey:" + key.ID
	tokenCache.SetTTL(token, principal, tokenTTL)
	c.Set("principal", principal)
}

func archiveFilter(c jug.Context) (*archive.Archive, bool) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	archiveId := c.Query(queryArchiveId)
	if archiveId == "" {
		handleError(c, ec.InvalidArgument)
		return nil, false
	}
	arch, err := archive.FindOne(c, b.Name, key, archiveId)
	if err != nil {
		handleError(c, err)
		return nil, false
	}
	return arch, true
}
