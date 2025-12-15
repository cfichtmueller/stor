// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"
	"log/slog"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

func handleRpcLogoutSession(c *srv.Context) *srv.Response {
	principal := contextMustGetPrincipal(c)
	userId, err := user.IdFromUrn(principal)
	if err != nil {
		return responseFromError(err)
	}
	sessionId := c.Query("session")
	if sessionId == "" {
		return responseFromError(errors.New("session ID is required"))
	}
	s, err := session.Get(c, sessionId)
	if err != nil {
		return responseFromError(err)
	}
	if s.User != userId {
		return responseFromError(errors.New("session not found"))
	}
	if err := session.Delete(c, sessionId); err != nil {
		slog.Error("failed to logout session", "error", err)
		return nodeResponse(e.Div(e.Class("col-span-9 text-red-500"), e.Raw("Failed to logout session")))
	}
	return nodeResponse(e.Div(e.Class("col-span-9 text-green-500"), e.Raw("Session logged out")))
}
