// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"strings"
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/archive"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/nonce"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/util"
)

var (
	tokenCache = util.NewCache()
	tokenTTL   = time.Minute
)

func bucketFilter(c jug.Context) {
	name := c.Param(paramBucketName)
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

func mustGetBucket(c jug.Context) (*bucket.Bucket, bool) {
	name := c.Param(paramBucketName)
	if len(name) < 3 {
		handleError(c, ec.NoSuchBucket)
		return nil, false
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		handleError(c, err)
		return nil, false
	}
	return b, true
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

func mustGetObject(c jug.Context, b *bucket.Bucket) (*object.Object, bool) {
	key := contextGetObjectKey(c)

	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil {
		handleError(c, err)
		return nil, false
	}

	return o, true
}

func authenticatedFilter(c jug.Context) {
	principal, ok := authenticateApiKey(c)
	if !ok {
		handleError(c, ec.Unauthorized)
		c.Abort()
		return
	}
	c.Set("principal", principal)
}

func authenticateApiKey(c jug.Context) (string, bool) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", false
	}
	token := auth[7:]
	if p, ok := tokenCache.Get(token); ok {
		return p.(string), true
	}
	key, err := apikey.Authenticate(c, token)
	if err != nil {
		handleError(c, err)
		return "", false
	}
	principal := "apikey:" + key.ID
	tokenCache.SetTTL(token, principal, tokenTTL)
	return principal, true
}

func mustAuthenticateApiKey(c jug.Context) bool {
	if _, ok := authenticateApiKey(c); !ok {
		handleError(c, ec.Unauthorized)
		return false
	}
	return true
}

func authenticateNonce(c jug.Context) (bool, error) {
	id := c.Query("nonce")
	if id == "" {
		return false, nil
	}
	n, err := nonce.GetAndInvalidate(c, id)
	if err != nil {
		if errors.Is(err, nonce.ErrNotFound) {
			return false, ec.Unauthorized
		}
		return false, err
	}
	if n.Bucket != c.Param(paramBucketName) || n.Key != contextGetObjectKey(c) {
		return false, ec.Unauthorized
	}
	return true, nil
}

func mustAuthenticateNonce(c jug.Context) bool {
	ok, err := authenticateNonce(c)
	if err != nil {
		handleError(c, err)
		return false
	}
	if !ok {
		handleError(c, ec.Unauthorized)
		return false
	}
	return true
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
