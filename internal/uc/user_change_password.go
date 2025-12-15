// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/ec"
)

type UserChangePasswordCommand struct {
	CurrentPassword string
	NewPassword     string
}

type UserChangePasswordResult struct {
	Message string
}

func UserChangePassword(ctx context.Context, u *user.User, cmd *UserChangePasswordCommand) (*UserChangePasswordResult, error) {
	if !u.PasswordMatches(cmd.CurrentPassword) {
		return &UserChangePasswordResult{
			Message: "Invalid current password",
		}, ec.InvalidArgument
	}
	if len(cmd.NewPassword) < 8 {
		return &UserChangePasswordResult{
			Message: "New password is too short",
		}, ec.InvalidArgument
	}

	if err := u.SetPassword(cmd.NewPassword); err != nil {
		return &UserChangePasswordResult{
			Message: "Internal error",
		}, err
	}

	if err := user.Update(ctx, u); err != nil {
		return &UserChangePasswordResult{
			Message: "Internal error",
		}, err
	}

	return &UserChangePasswordResult{}, nil
}
