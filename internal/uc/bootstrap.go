// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type BootstrapCommand struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

func (c BootstrapCommand) Validate() error {
	v := srv.RequireNotEmpty("email", c.Email, nil)
	v = srv.RequireNotEmpty("password", c.Password, v)
	v = srv.Require("password", "password_mismatch", "Passwords do not match", c.Password == c.PasswordConfirmation, v)
	return srv.Validate(v)
}

func Bootstrap(ctx context.Context, cmd BootstrapCommand) error {
	if _, err := user.Create(ctx, user.CreateCommand{
		Email:    cmd.Email,
		Password: cmd.Password,
	}); err != nil {
		return err
	}
	return nil
}
