// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

const (
	app_sidebar_active_dashboard = "dashboard"
	app_sidebar_active_buckets   = "buckets"
	app_sidebar_active_admin     = "admin"
)

func appSidebar(active string) e.Node {
	return Sidebar(
		SidebarItem{Title: "Dashboard", Link: dashboardLink, Active: active == app_sidebar_active_dashboard, Icon: IconGauge},
		SidebarItem{Title: "Buckets", Link: bucketsLink, Active: active == app_sidebar_active_buckets, Icon: IconArchive},
		SidebarItem{Title: "Admin", Link: adminLink, Active: active == app_sidebar_active_admin, Icon: IconCog},
	)
}
