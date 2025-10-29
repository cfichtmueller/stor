// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"
	"fmt"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/ui"
)

var (
	ErrLoginRequired = fmt.Errorf("login required")
)

func bucketFilter(c *srv.Context, next srv.Handler) *srv.Response {
	name := c.PathValue("bucketName")
	if len(name) < 3 {
		return nodeResponse(ui.NotFoundPage())
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		return responseFromError(err)
	}
	c.Set("bucket", b)
	return next(c)
}

func withBucketFromQuery(c *srv.Context, next srv.Handler) *srv.Response {
	name := c.Query("bucket")
	if name == "" {
		return next(c)
	}
	b, err := bucket.FindOne(c, name)
	if err != nil {
		return responseFromError(err)
	}
	contextSetBucket(c, b)
	return next(c)
}

func apiKeyFilter(c *srv.Context, next srv.Handler) *srv.Response {
	keyId := c.Query("key")
	if keyId == "" {
		return next(c)
	}
	key, err := apikey.Get(c, keyId)
	if err != nil {
		return responseFromError(err)
	}
	c.Set("apiKey", key)
	return next(c)
}

func authenticatedFilter(c *srv.Context, next srv.Handler) *srv.Response {
	authCookie, err := c.Cookie("stor_auth")
	if err != nil {
		return hxRedirect(c, "/login")
	}
	userId, err := authenticateSession(c, authCookie)
	if err != nil {
		if errors.Is(err, ErrLoginRequired) {
			return hxRedirect(c, "/login")
		}
		return responseFromError(err)
	}
	contextSetPrincipal(c, user.Urn(userId))
	return next(c)
}

// authenticates a session and returns the user's id if successful
func authenticateSession(c *srv.Context, sessionId string) (string, error) {
	s, err := session.Get(c, sessionId)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return "", ErrLoginRequired
		}
		return "", err
	}
	if s.IsExpired() {
		return "", ErrLoginRequired
	}
	return s.User, nil
}
