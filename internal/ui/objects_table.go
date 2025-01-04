// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func ObjectsTable(objects []ObjectData) e.Node {
	return Table(
		TableHeader(
			TableHead("", e.Text("Key")),
			TableHead("", e.Text("Size")),
		),
		TableBody(
			e.Mapf(objects, func(o ObjectData) e.Node {
				return TableRow(
					TableCell(
						e.A(
							e.Href(o.Href),
							e.Text(o.Key),
						),
					),
					TableCell(
						e.Iff(o.Size > 0, e.F(bytesText, o.Size)),
					),
				)
			}),
		),
	)
}
