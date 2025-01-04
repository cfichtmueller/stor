// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func LoginPage() e.Node {
	return AuthLayout(
		e.Div(
			e.Class("flex flex-col space-y-2 text-center"),
			e.H1(
				e.Class("text-2xl font-semibold tracking-tight"),
				e.Text("Log In"),
			),
		),
		LoginForm(LoginFormData{}),
	)
}
