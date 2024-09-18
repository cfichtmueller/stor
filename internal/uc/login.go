// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type LoginCommand struct {
	IpAddress string
	Email     string
	Password  string
}

// Login logs in a user. Returns the user and a session token on success
func Login(ctx context.Context, cmd LoginCommand) (*user.User, string, error) {
	u, err := user.Login(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, "", err
	}

	s, err := session.Create(ctx, u.ID, cmd.IpAddress)
	if err != nil {
		return nil, "", err
	}

	return u, s.ID, nil
}
