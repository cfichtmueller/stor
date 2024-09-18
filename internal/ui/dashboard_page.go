// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"
)

type DashboardPageModel struct {
	PageHeader pageHeaderModel
	AppSidebar SidebarModel
	Metrics    dashboardMetricsModel
}

func RenderDashboardPage(w io.Writer, d DashboardData) error {
	return renderTemplate(w, "DashboardPage", DashboardPageModel{
		PageHeader: pageHeaderModel{
			Title: "Dashboard",
		},
		AppSidebar: appSidebarModel(app_sidebar_active_dashboard),
		Metrics:    newDashboardMetricsModel(d),
	})
}
