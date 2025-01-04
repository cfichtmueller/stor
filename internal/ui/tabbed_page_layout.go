// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
)

func TabbedPageLayout(sidebar, pageHeader, navTabs e.Node, children ...e.Node) e.Node {
	return LoggedInLayout(
		sidebar,
		e.Div(
			e.Class("w-full min-h-full flex flex-col px-2"),
			pageHeader,
			e.Div(
				e.Class("flex bg-neutral-50 flex-grow mb-4"),
				e.Div(
					e.Class("flex flex-col w-full border-r border border-l-0 rounded-md bg-white"),
					e.Div(
						e.Class("flex gap-2 w-full py-2 ps-2 text-sm"),
						navTabs,
					),
					e.Div(
						append(children, e.Class("p-2"))...,
					),
				),
			),
		),
	)
}
