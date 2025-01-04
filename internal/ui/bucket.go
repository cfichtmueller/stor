// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
)

func CreateBucketDialog() e.Node {
	return e.Form(
		e.HXPost("/r/bucket"),
		e.HXSwap(e.HXSwapDelete),
		e.Id("create-bucket-dialog"),
		DialogBackdrop(),
		e.Div(
			e.Role("dialog"),
			e.Class("fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg"),
			e.TabIndex(-1),
			e.StyleAttr("pointer-events: auto;"),
			DialogTitle(e.Raw("Create Bucket")),
			e.Div(
				e.Class("grid gap-4 py-4"),
				e.Div(
					e.Class("grid grid-cols-4 items-center gap-4"),
					e.Label(
						e.For("bucketName"),
						e.Class(cn(cnLabel, "text-right")),
						e.Raw("Name"),
					),
					e.Input(
						e.Id("bucketName"),
						e.Name("name"),
						e.Class(cn(cnInput, "col-span-3")),
						e.Required(),
					),
				),
			),
			e.Div(
				e.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"),
				e.Button(
					e.Class(cn(btn, btnPrimary)),
					e.Type("submit"),
					e.Raw("Create"),
				),
			),
			e.Button(
				e.Type("button"),
				e.Class(cnDIalogCloseButton),
				e.Attr("data-remove", "create-bucket-dialog"),
				IconDialogClose,
				e.Span(srOnly(), e.Raw("Close")),
			),
		),
	)
}
