// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type UsersPageData struct {
	Users []*user.User
}

func UsersPage(d *UsersPageData) e.Node {
	return AdminPageLayout(
		admin_tab_active_users,
		UsersTable(d.Users),
	)
}
