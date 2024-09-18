// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type BootstrapCommand struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

func (c BootstrapCommand) Validate() error {
	return jug.NewValidator().
		RequireNotEmpty(c.Email, "Email is missing").
		RequireNotEmpty(c.Password, "Password is missing").
		Require(c.Password == c.PasswordConfirmation, "Passwords do not match").
		Validate()
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
