// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

type BootstrapFormData struct {
	Email        string
	ErrorMessage string
}

func BootstrapForm(f *BootstrapFormData) e.Node {
	return e.Form(
		e.HXPost("/bootstrap"),
		e.HXSwap("outerHTML"),
		e.Div(
			e.Class("grid gap-6"),
			FormLabel("email", "Email"),
			e.Input(
				e.Type("email"),
				e.Id("email"),
				e.Name("email"),
				e.Placeholder("Email"),
				e.Class(cnInput),
				e.Required(),
				e.AutoComplete("username"),
				e.Value(f.Email),
			),
			FormLabel("password", "Password"),
			e.Input(
				e.Type("password"),
				e.Id("password"),
				e.Name("password"),
				e.Placeholder("Password"),
				e.Class(cnInput),
				e.Required(),
				e.AutoComplete("new-password"),
			),
			FormLabel("passwordConfirmation", "Confirm Password"),
			e.Input(
				e.Type("password"),
				e.Id("password"),
				e.Name("passwordConfirmation"),
				e.Placeholder("Repeat Password"),
				e.Class(cnInput),
				e.Required(),
				e.AutoComplete("new-password"),
			),
			e.If(f.ErrorMessage != "", e.Div(
				e.Class("text-red-400"),
				e.Text(f.ErrorMessage),
			)),
			e.Button(
				e.Type("submit"),
				e.Class(cn(btn, btnPrimary)),
				e.Text("Create admin account"),
			),
		),
	)
}
