// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func ListPageLayout(title string, sidebar e.Node, children ...e.Node) e.Node {
	return LoggedInLayout(
		sidebar,
		e.Div(
			e.Class("w-full flex-grow px-2"),
			PageHeader(title, ""),
			e.Div(
				append(children, e.Class("rounded-lg bg-white border flex-grow flex flex-col p-4 w-full"))...,
			),
		),
	)
}
