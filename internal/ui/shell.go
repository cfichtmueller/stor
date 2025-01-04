// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func Shell(title string, children ...e.Node) e.Node {
	t := "STOR"
	if title != "" {
		t = title + " - STOR"
	}
	return e.Group(
		e.Doctype(),
		e.Html(
			e.Lang("en"),
			e.Head(
				e.Title(e.Raw(t)),
				e.Meta(e.Charset("UTF-8")),
				e.Meta(e.Name("viewport"), e.Content("width=device-width, initial-scale=1.0")),
				e.Link(e.Rel("icon"), e.Href("/img/icon.png?v=1726228022"), e.Type("image/png")),
				e.Link(e.RelStylesheet(), e.Href("/css/style.css?v=1726228022")),
			),
			e.Body(
				append(
					children,
					e.Class("min-h-screen"),
					ToastContainer(),
					e.Script(e.Src("/js/htmx.min.js")),
					e.Script(e.Src("/js/lib.js")),
				)...,
			),
		),
	)
}
