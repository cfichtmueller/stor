// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func Table(children ...e.Node) e.Node {
	return e.Div(
		e.Class("relative w-full overflow-auto"),
		e.Table(
			append(children, e.Class("w-full caption-bottom text-sm"))...,
		),
	)
}

func TableHeader(children ...e.Node) e.Node {
	return e.Thead(
		e.Class("[&_tr]:border-b"),
		e.Tr(
			append(children, e.Class("border-b"))...,
		),
	)
}

func TableHead(className string, children ...e.Node) e.Node {
	return e.Th(
		append(children, e.Class(cn("h-10 px-2 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]", className)))...,
	)
}

func TableBody(children ...e.Node) e.Node {
	return e.Tbody(append(children, e.Class("[&_tr:last-child]:border-0"))...)
}

func TableRow(children ...e.Node) e.Node {
	return e.Tr(append(children, e.Class("border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"))...)
}

func TableCell(children ...e.Node) e.Node {
	return TableCellC("", children...)
}

func TableCellC(className string, children ...e.Node) e.Node {
	return e.Td(
		append(children, e.Class(cn("p-2 align-middle [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]", className)))...,
	)
}
