// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

type pageHeaderModel struct {
	Title     string
	CloseLink string
}

func PageHeader(title, closeLink string) e.Node {
	return e.Div(
		e.Class("flex justify-between py-4"),
		e.Div(
			e.Class("flex items-center gap-x-2"),
			e.If(closeLink != "", e.A(
				e.Href(closeLink),
				IconX,
			)),
			e.Div(
				e.Class("flex text-sm items-center gap-x-2"),
				e.Raw(title),
			),
		),
		e.Div(),
	)
}
