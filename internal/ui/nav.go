// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

type NavLink struct {
	Active bool
	Link   string
	Title  string
	Icon   e.Node
}

type NavTabsModel struct {
	Tabs []*NavLink
}

func NewNavTabsModel() *NavTabsModel {
	return &NavTabsModel{
		Tabs: make([]*NavLink, 0),
	}
}

type SidebarModel struct {
	Items []NavLink
}

func NavTabs(links ...*NavLink) e.Node {
	return e.Div(
		e.Class("flex gap-2 border-b w-full py-2 ps-2 text-sm"),
		e.HXBoost(),
		e.Mapf(links, navLinkNode),
	)
}

func navLinkNode(l *NavLink) e.Node {
	if l.Active {
		return e.Div(
			e.Class("flex text-sm items-center rounded-sm p-1 font-medium bg-neutral-100 cursor-pointer"),
			l.Icon,
			e.Raw(l.Title),
		)
	}
	return e.A(
		e.Class("flex text-sm items-center rounded-sm p-1"),
		e.Href(l.Link),
		l.Icon,
		e.Text(l.Title),
	)
}
