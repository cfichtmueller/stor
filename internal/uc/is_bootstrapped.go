// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uc

import (
	"context"

	"github.com/cfichtmueller/stor/internal/domain/user"
)

var isBootstrappedInitialized = false
var isBootstrapped = false

func IsBootstrapped(ctx context.Context) (bool, error) {
	if isBootstrappedInitialized {
		return isBootstrapped, nil
	}
	users, err := user.List(ctx)
	if err != nil {
		return false, err
	}
	isBootstrapped = len(users) > 0
	isBootstrappedInitialized = isBootstrapped
	return isBootstrapped, nil
}
