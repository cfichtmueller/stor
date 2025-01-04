// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
)

func ApiKeysTable(keys []*apikey.ApiKey) e.Node {
	return Table(
		TableHeader(
			TableHead("", e.Text("Key")),
			TableHead("", e.Text("Description")),
			TableHead("text-right", e.Text("Created At")),
			TableHead("text-right", e.Text("Created By")),
			TableHead("text-right", e.Text("Expires")),
			TableHead("flex justify-end", e.Button(
				e.Class(cn(btn, btnPrimary)),
				e.HXGet("/c/create-api-key-dialog"),
				e.HXTarget("body"),
				e.HXSwap("beforeend"),
				IconPlus,
				e.Span(
					srOnly(),
					e.Text("Create"),
				),
			)),
		),
		TableBody(
			e.Mapf(keys, ApiKeysTableRow),
		),
	)
}

func ApiKeysTableRow(key *apikey.ApiKey) e.Node {
	return TableRow(
		TableCellC("p-2 w-1 align-middle whitespace-nowrap", e.Text(key.Prefix)),
		TableCellC("cursor-pointer",
			e.HXGet("/c/api-key-sheet?key="+key.ID),
			e.HXTarget("body"),
			e.HXSwap("beforeend"),
			e.Text(key.Description),
		),
		TableCellC("text-right", e.Text(formatDateTime(key.CreatedAt))),
		TableCellC("text-right", e.Text(key.CreatedBy)),
		TableCellC("text-right", e.Text(formatDateTime(key.ExpiresAt))),
		TableCellC("p-2 flex w-full justify-end align-middle"),
	)
}
