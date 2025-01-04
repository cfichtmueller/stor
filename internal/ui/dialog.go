// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

type dialogModel struct {
	ID           string
	Title        string
	Description  string
	CancelTitle  string
	ConfirmTitle string
	Destructive  bool
	HxDelete     string
	HxPost       string
	HxTarget     string
	HxSwap       string
}

func Dialog(id string, children ...e.Node) e.Node {
	return e.Div(
		e.Id(id),
		DialogBackdrop(),
		DialogContent(
			e.Button(
				e.Type("button"),
				e.Class(cnDIalogCloseButton),
				e.Attr("data-remove", id),
				IconClose,
				e.Span(srOnly(), e.Text("Close")),
			),
			e.Group(children...),
		),
	)
}

func DialogContent(children ...e.Node) e.Node {
	return e.Div(
		e.Role("dialog"),
		e.Class("fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg"),
		e.TabIndex(-1),
		e.StyleAttr("pointer-events: auto;"),
		e.Group(children...),
	)
}

func DialogBackdrop() e.Node {
	return e.Div(
		e.Class("fixed inset-0 z-50 bg-black/80 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"),
		e.StyleAttr("pointer-events: auto;"),
		e.Attr("data-aria-hidden", "true"),
		e.AriaHidden(),
	)
}

func DialogTitle(children ...e.Node) e.Node {
	return e.H2(
		append(children, e.Class("text-lg font-semibold leading-none tracking-tight"))...,
	)
}

type ConfirmationDialogData struct {
	ID           string
	Title        string
	Description  string
	CancelTitle  string
	Destructive  bool
	HxPost       string
	HxDelete     string
	HxTarget     string
	HxSwap       string
	ConfirmTitle string
}

func ConfirmationDialog(d *ConfirmationDialogData) e.Node {
	btnClass := btnPrimary
	if d.Destructive {
		btnClass = btnDanger
	}
	return Dialog(
		d.ID,
		DialogTitle(e.Text(d.Title)),
		e.Div(e.Text(d.Description)),
		e.Div(
			e.Class("flex justify-between"),
			e.Button(
				e.Class(cn(btn, btnSecondary)),
				e.Attr("data-remove", d.ID),
				e.Text(d.CancelTitle),
			),
			e.Button(
				e.Class(cn(btn, btnClass)),
				e.If(d.HxPost != "", e.HXPost(d.HxPost)),
				e.If(d.HxDelete != "", e.HXDelete(d.HxDelete)),
				e.If(d.HxTarget != "", e.HXTarget(d.HxTarget)),
				e.If(d.HxSwap != "", e.HXSwap(d.HxSwap)),
				e.Text(d.ConfirmTitle),
			),
		),
	)
}
