// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

type Breadcrumb struct {
	Separator bool
	Link      string
	Title     string
}

func Breadcrumbs(crumbs ...Breadcrumb) e.Node {
	return e.Nav(
		ariaLabel("breadcrumb"),
		e.Ol(
			e.Class("flex flex-wrap items-center break-words text-sm text-muted-foreground sm:gap-1"),
			e.Mapf(crumbs, func(c Breadcrumb) e.Node {
				if c.Separator {
					return e.Li(
						ariaRole("presentation"),
						e.AriaHidden(),
						e.Class("[&amp;>svg]:size-3.5"),
						e.Raw(`<svg width="15" height="15" viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M4.10876 14L9.46582 1H10.8178L5.46074 14H4.10876Z" fill="currentColor" fill-rule="evenodd" clip-rule="evenodd"></path></svg>`),
					)
				} else if c.Link != "" {
					return e.Li(
						e.HXBoost(),
						e.Class("inline-flex items-center"),
						e.A(
							e.Class("transition-colors hover:text-foreground"),
							e.Href(c.Link),
							e.Text(c.Title),
						),
					)
				}
				return e.Li(
					e.Class("inline-flex items-center"),
					e.Span(
						e.Role("link"),
						ariaDisabled(),
						ariaCurrent("page"),
						e.Class("font-normal text-foreground"),
						e.Text(c.Title),
					),
				)
			}),
		),
	)
}
