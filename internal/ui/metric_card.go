// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func MetricCard(title string, icon e.Node, value, lead string) e.Node {
	return e.Div(
		e.Class("rounded-xl border bg-card text-card-foreground shadow"),
		e.Div(
			e.Class("p-6 flex flex-row items-center justify-between space-y-0 pb-2"),
			e.H3(
				e.Class("tracking-tight text-sm font-medium"),
				e.Raw(title),
			),
			icon,
		),
		e.Div(
			e.Class("p-6 pt-0"),
			e.Div(
				e.Class("text-2xl font-bold"),
				e.Raw(value),
			),
			e.If(lead != "", e.P(
				e.Class("text-xs text-muted-foreground"),
				e.Raw(lead),
			)),
		),
	)
}
