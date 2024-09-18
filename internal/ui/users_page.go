// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/user"
)

type UsersPageData struct {
	Users []*user.User
}

type usersPageModel struct {
	Layout adminPageModel
	Users  []*user.User
}

func RenderUsersPage(w io.Writer, d UsersPageData) error {
	return renderTemplate(w, "UsersPage", usersPageModel{
		Layout: newAdminPageModel(admin_tab_active_users),
		Users:  d.Users,
	})
}
