// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func AuthLayout(children ...e.Node) e.Node {
	return Shell(
		"Bootstrap",
		e.Div(
			e.Class("container relative flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0 min-h-[100vh]"),
			e.Div(
				e.Class("relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex"),
				e.Div(
					e.Class("absolute inset-0 bg-zinc-900"),
				),
				e.Div(
					e.Class("relative z-20 flex items-center text-lg font-medium"),
					e.Img(
						e.Class("mr-2 h-6 w-6"),
						e.Src("/img/icon.png?v=1726228022"),
					),
					e.Text("STOR"),
				),
				e.Div(
					e.Class("relative z-20 m-auto"),
					e.Img(
						e.Class("h-96"),
						e.Src("/img/icon.png?v=1726228022"),
					),
				),
			),
			e.Div(
				e.Class("lg:p-8"),
				e.Div(
					append(children, e.Class("mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]"))...,
				),
			),
		),
	)
}
