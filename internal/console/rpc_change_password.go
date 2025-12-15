// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"log/slog"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleRpcChangePassword(c *srv.Context) *srv.Response {
	principal := contextMustGetPrincipal(c)
	id, err := user.IdFromUrn(principal)
	if err != nil {
		return responseFromError(err)
	}
	u, err := user.Get(c, id)
	if err != nil {
		return responseFromError(err)
	}
	values := c.FormValues()
	currentPassword := values.Get("currentPassword")
	newPassword := values.Get("newPassword")

	m, err := uc.UserChangePassword(c, u, &uc.UserChangePasswordCommand{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	})

	if err != nil {
		slog.Error("failed to change user password", "error", err)
	}

	return nodeResponseWithShell(c, ui.ProfilePasswordSection(&ui.ProfilePasswordSectionData{
		Message: m.Message,
	}))
}
