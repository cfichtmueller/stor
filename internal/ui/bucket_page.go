// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func BucketPage(links *BucketLinks, activeTab string, breadcrumbs, pageTitle e.Node, children ...e.Node) e.Node {
	return LoggedInLayout(
		appSidebar(app_sidebar_active_buckets),
		breadcrumbs,
		pageTitle,
		e.Div(
			e.Class("flex flex-col w-full border rounded-md bg-white"),
			BucketNavTabs(links, activeTab),
			e.Div(
				e.Class("p-2"),
				e.Group(children...),
			),
		),
	)
}
