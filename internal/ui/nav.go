// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

type NavLink struct {
	Active bool
	Link   string
	Title  string
	Icon   string
}

type NavTabsModel struct {
	Tabs []NavLink
}

type SidebarModel struct {
	Items []NavLink
}
