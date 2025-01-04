// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

func UsersTable(users []*user.User) e.Node {
	return Table(
		TableHeader(
			TableHead("", e.Text("Email")),
			TableHead("", e.Text("Status")),
			TableHead("flex justify-end"),
		),
		TableBody(
			e.Mapf(users, UsersTableRow),
		),
	)
}

func UsersTableRow(u *user.User) e.Node {
	return TableRow(
		TableCellC("p-2 w-1 align-middle whitespace-nowrap", e.Text(u.Email)),
		TableCell(
			e.If(u.Enabled, e.Text("Enabled")),
			e.If(!u.Enabled, e.Text("Disabled")),
		),
		TableCellC("p-2 flex w-full justify-end align-middle"),
	)
}
