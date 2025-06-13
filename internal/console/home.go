// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/uc"
)

func handleHomePage(c *srv.Context) *srv.Response {
	bootstrapped, err := uc.IsBootstrapped(c)
	if err != nil {
		return srv.Respond().Error(err)
	}
	if !bootstrapped {
		return srv.Respond().Found("/bootstrap")
	}
	//TODO: redirect to /u when user is authenticated
	return srv.Respond().Found("/login")
}
