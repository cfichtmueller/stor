// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/uc"
)

func handleHomePage(c jug.Context) {
	bootstrapped, err := uc.IsBootstrapped(c)
	if err != nil {
		c.HandleError(err)
		return
	}
	if !bootstrapped {
		redirect(c, "/bootstrap")
		return
	}
	//TODO: redirect to /u when user is authenticated
	redirect(c, "/login")
}
