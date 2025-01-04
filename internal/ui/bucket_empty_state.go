// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func BucketEmptyState() e.Node {
	return e.Div(
		e.Class("flex flex-col justify-center items-center min-h-96"),
		e.Img(
			e.Src("/img/empty.png"),
			e.Class("h-32 mb-8"),
			e.AriaHidden(),
		),
		e.Div(e.Text("This bucket is empty")),
	)
}
