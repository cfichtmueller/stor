// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
)

type apiKeyModel struct {
	ID          string
	Prefix      string
	Description string
	CreatedAt   string
	CreatedBy   string
	ExpiresAt   string
}

func newApiKeyModel(k *apikey.ApiKey) apiKeyModel {
	return apiKeyModel{
		ID:          k.ID,
		Prefix:      k.Prefix,
		Description: k.Description,
		CreatedAt:   formatDateTime(k.CreatedAt),
		CreatedBy:   k.CreatedBy,
		ExpiresAt:   formatDateTime(k.ExpiresAt),
	}
}

func ApiKeyCreatedDialog(k *apikey.ApiKey, plain string) e.Node {
	return e.Div(
		e.Id("api-key-created-dialog"),
		DialogBackdrop(),
		e.Div(
			e.Role("dialog"),
			e.Class("fixed left-[50%] top-[50%] z-50 grid w-full max-w-3xl translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg"),
			e.TabIndex(-1),
			e.StyleAttr("pointer-events: auto;"),
			DialogTitle(e.Raw("API Key Created")),
			e.Div(
				e.Class("grid gap-4 py-4"),
				e.Div(e.Raw("Your API key has been created. Please copy your key since you won't be able to access it again")),
				e.Div(e.Raw("Your API key")),
				e.Div(
					e.Class("border p-2 flex items-center"),
					e.Div(
						e.Class("flex-grow"),
						e.Raw(plain),
					),
					e.Div(
						e.Attr("data-clipboard", plain),
						IconCopy,
					),
				),
			),
			e.Div(
				e.Class("flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2"),
				e.Button(
					e.Class(cn(btn, btnPrimary)),
					e.Attr("data-remove", "api-key-created-dialog"),
					e.Raw("OK"),
				),
			),
			e.Button(
				e.Type("button"),
				e.Class(cnDIalogCloseButton),
				e.Attr("data-remove", "api-key-created-dialog"),
				IconDialogClose,
				e.Span(
					srOnly(),
					e.Raw("Close"),
				),
			),
		),
	)
}

func ApiKeySheet(key *apikey.ApiKey) e.Node {
	return Sheet(
		sheetModel{
			Title: "API Key",
			Lead:  "Make changes to the API key here.",
		},
		e.Div(
			e.Class("flex justify-end"),
			e.Button(
				e.Class(cn(btn, btnDanger)),
				e.HXGet("/c/api-key-delete-dialog?key="+key.ID),
				e.HXTarget("body"),
				e.HXSwap(e.HXSwapBeforeend),
				e.Raw("Delete"),
			),
		),
	)
}

func CreateApiKeyDialog() e.Node {
	return e.Form(
		e.HXPost("/r/api-key"),
		e.HXSwap("delete"),
		e.Id("create-api-key-dialog"),
		DialogBackdrop(),
		DialogContent(
			DialogTitle(e.Text("Create API Key")),
			e.Div(
				e.Class("grid gap-4 py-4"),
				e.Div(
					e.Class("grid grid-cols-4 items-center gap-4"),
					e.Label(
						e.For("description"),
						e.Class(cn(cnLabel, "text-right")),
						e.Text("Description"),
					),
					e.Input(
						e.Class(cn(cnInput, "col-span-3")),
						e.Name("description"),
						e.Required(),
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
					e.Attr("data-remove", "create-api-key-dialog"),
					IconDialogClose,
					e.Span(srOnly(), e.Raw("Close")),
				),
			),
		),
	)
}

func DeleteApiKeyDialog(key *apikey.ApiKey) e.Node {
	return ConfirmationDialog(&ConfirmationDialogData{
		ID:           "delete-api-key-dialog",
		Title:        "Delete API Key",
		Description:  "U sure u want to delete key " + key.Description + "?",
		ConfirmTitle: "Delete",
		CancelTitle:  "Cancel",
		Destructive:  true,
		HxDelete:     "/r/api-key?key=" + key.ID,
	})
}
