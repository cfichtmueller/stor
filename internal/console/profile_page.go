// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleProfilePage(c *srv.Context) *srv.Response {
	principal := contextMustGetPrincipal(c)
	id, err := user.IdFromUrn(principal)
	if err != nil {
		return responseFromError(err)
	}
	u, err := user.Get(c, id)
	if err != nil {
		return responseFromError(err)
	}
	return nodeResponseWithShell(c, ui.ProfilePage(&ui.ProfilePageData{
		User: u,
	}))
}
