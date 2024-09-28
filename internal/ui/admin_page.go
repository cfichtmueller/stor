// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

var (
	admin_tab_active_api_keys = "apiKeys"
	admin_tab_active_users    = "users"
)

type adminPageModel struct {
	PageHeader pageHeaderModel
	AppSidebar SidebarModel
	NavTabs    *NavTabsModel
}

func newAdminPageModel(activeTab string) adminPageModel {
	return adminPageModel{
		PageHeader: pageHeaderModel{
			Title: "Admin",
		},
		AppSidebar: appSidebarModel(app_sidebar_active_admin),
		NavTabs: &NavTabsModel{
			Tabs: []*NavLink{
				{
					Link:   usersLink,
					Active: activeTab == admin_tab_active_users,
					Title:  "Users",
					Icon:   "users-round",
				},
				{
					Link:   apiKeysLink,
					Active: activeTab == admin_tab_active_api_keys,
					Title:  "API Keys",
					Icon:   "key-round",
				},
			},
		},
	}
}
