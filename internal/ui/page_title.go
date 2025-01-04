// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func PageTitle(title string) e.Node {
	return e.Div(
		e.Class("text-2xl font-medium pb-4"),
		e.Text(title),
	)
}
