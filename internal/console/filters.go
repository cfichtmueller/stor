// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"
	"fmt"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/ui"
)

var (
	ErrLoginRequired = fmt.Errorf("login required")
)

func bucketFilter(c jug.Context) {
	name := c.Param("bucketName")
	if len(name) < 3 {
		must("render not found page", c, ui.RenderNotFoundPage(c.Writer()))
		c.Abort()
		return
	}
	b, err := bucket.FindOne(c, name)
	if !must("find bucket", c, err) {
		c.Abort()
		return
	}
	c.Set("bucket", b)
}

func apiKeyFilter(c jug.Context) {
	keyId := c.Query("key")
	if keyId == "" {
		c.Next()
		return
	}
	key, err := apikey.Get(c, keyId)
	if err != nil {
		c.HandleError(err)
		c.Abort()
		return
	}
	c.Set("apiKey", key)
}

func authenticatedFilter(c jug.Context) {
	authCookie, authCookieExists := c.Cookie("stor_auth")
	if authCookieExists {
		if err := authenticateSession(c, authCookie); err != nil {
			if errors.Is(err, ErrLoginRequired) {
				hxRedirect(c, "/login")
				c.Abort()
				return
			}
			c.HandleError(err)
			c.Abort()
			return
		}
		return
	}
	hxRedirect(c, "/login")
	c.Abort()
}

func authenticateSession(c jug.Context, sessionId string) error {
	s, err := session.Get(c, sessionId)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return ErrLoginRequired
		}
		return err
	}
	if s.IsExpired() {
		return ErrLoginRequired
	}
	return nil
}
