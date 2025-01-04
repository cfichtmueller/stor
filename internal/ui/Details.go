// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func Detail(title, value string) e.Node {
	return e.Div(
		e.Class("flex items-center py-1 font-light"),
		e.Div(
			e.Class("text-neutral-400 text-sm whitespace-nowrap w-1/3"),
			e.Text(title),
		),
		e.Div(
			e.Class("text-sm"),
			e.Text(value),
		),
	)
}

func Details(title string, details ...e.Node) e.Node {
	return e.Div(
		e.If(title != "", e.Div(
			e.Class("pb-2 font-semibold text-sm text-neutral-500"),
			e.Text(title),
		)),
		e.Mapf(details, identity),
	)
}
