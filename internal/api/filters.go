// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"strings"
	"time"

	"github.com/cfichtmueller/srv"
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

func bucketFilter(c *srv.Context, next srv.Handler) *srv.Response {
	name := c.PathValue(paramBucketName)
	if len(name) < 3 {
		return responseFromError(ec.NoSuchBucket)
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		return responseFromError(err)
	}
	c.Set("bucket", b)
	return next(c)
}

func mustGetBucket(c *srv.Context) (*bucket.Bucket, *srv.Response) {
	name := c.PathValue(paramBucketName)
	if len(name) < 3 {
		return nil, responseFromError(ec.NoSuchBucket)
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		return nil, responseFromError(err)
	}
	return b, nil
}

func objectFilter(c *srv.Context) (*object.Object, *srv.Response) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)

	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil {
		return nil, responseFromError(err)
	}

	return o, nil
}

func mustGetObject(c *srv.Context, b *bucket.Bucket) (*object.Object, *srv.Response) {
	key := contextGetObjectKey(c)

	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil {
		return nil, responseFromError(err)
	}

	return o, nil
}

func authenticatedFilter(c *srv.Context, next srv.Handler) *srv.Response {
	principal, ok, err := authenticateApiKey(c)
	if err != nil {
		return responseFromError(err)
	}
	if !ok {
		return responseFromError(ec.Unauthorized)
	}
	c.Set("principal", principal)
	return next(c)
}

func authenticateApiKey(c *srv.Context) (string, bool, error) {
	auth := c.Header("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", false, nil
	}
	token := auth[7:]
	if p, ok := tokenCache.Get(token); ok {
		return p.(string), true, nil
	}
	key, err := apikey.Authenticate(c, token)
	if err != nil {
		return "", false, err
	}
	principal := "apikey:" + key.ID
	tokenCache.SetTTL(token, principal, tokenTTL)
	return principal, true, nil
}

func mustAuthenticateApiKey(c *srv.Context) *srv.Response {
	_, ok, err := authenticateApiKey(c)
	if err != nil {
		return responseFromError(err)
	}
	if !ok {
		return responseFromError(ec.Unauthorized)
	}
	return nil
}

func authenticateNonce(c *srv.Context) (bool, error) {
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
	if n.Bucket != c.PathValue(paramBucketName) || n.Key != contextGetObjectKey(c) {
		return false, ec.Unauthorized
	}
	return true, nil
}

func mustAuthenticateNonce(c *srv.Context) *srv.Response {
	ok, err := authenticateNonce(c)
	if err != nil {
		return responseFromError(err)
	}
	if !ok {
		return responseFromError(ec.Unauthorized)
	}
	return nil
}

func archiveFilter(c *srv.Context) (*archive.Archive, *srv.Response) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	archiveId := c.Query(queryArchiveId)
	if archiveId == "" {
		return nil, responseFromError(ec.InvalidArgument)
	}
	arch, err := archive.FindOne(c, b.Name, key, archiveId)
	if err != nil {
		return nil, responseFromError(err)
	}
	return arch, nil
}
