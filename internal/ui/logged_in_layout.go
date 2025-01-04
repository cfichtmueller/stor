// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func LoggedInLayout(sidebar e.Node, children ...e.Node) e.Node {
	return e.Div(
		e.Class("flex flex-row min-h-screen w-full bg-neutral-100"),
		e.Div(
			e.Class("flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0 w-full max-w-screen-2xl"),
			sidebar,
			e.Div(append(children, e.Class("w-full min-h-full flex flex-col px-2 pt-4"))...),
		),
	)
}
