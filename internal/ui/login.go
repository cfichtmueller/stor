// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "io"

func RenderBootstrapPage(w io.Writer) error {
	return renderTemplate(w, "BootstrapPage", BootstrapFormModel{})
}

type BootstrapFormModel struct {
	Email        string
	ErrorMessage string
}

func RenderBootstrapForm(w io.Writer, model BootstrapFormModel) error {
	return renderTemplate(w, "BootstrapForm", model)
}

func RenderLoginPage(w io.Writer) error {
	return renderTemplate(w, "LoginPage", LoginFormModel{})
}

type LoginFormModel struct {
	Email        string
	ErrorMessage string
}

func RenderLoginForm(w io.Writer, model LoginFormModel) error {
	return renderTemplate(w, "LoginForm", model)
}
