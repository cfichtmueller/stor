// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

var (
	admin_tab_active_api_keys = "apiKeys"
	admin_tab_active_users    = "users"
)

func AdminPageLayout(active string, children ...e.Node) e.Node {
	return TabbedPageLayout(
		appSidebar(app_sidebar_active_admin),
		PageHeader("Admin", ""),
		NavTabs(
			&NavLink{
				Link:   usersLink,
				Active: active == admin_tab_active_users,
				Title:  "Users",
				Icon:   IconUsersRound,
			},
			&NavLink{
				Link:   apiKeysLink,
				Active: active == admin_tab_active_api_keys,
				Title:  "API Keys",
				Icon:   IconKeyRound,
			},
		),
		children...,
	)
}
