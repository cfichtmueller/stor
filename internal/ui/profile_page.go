// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type ProfilePageData struct {
	User *user.User
}

func ProfilePage(d *ProfilePageData) e.Node {
	return LoggedInLayout(
		appSidebar(app_sidebar_active_profile),
		e.Div(
			e.Class("w-full flex-grow px-2"),
			PageHeader("Profile", ""),
			e.Div(
				e.Class("rounded-lg bg-white border flex flex-col gap-y-4 p-4 w-full"),
				ProfilePasswordSection(&ProfilePasswordSectionData{}),
			),
		),
	)
}

type ProfilePasswordSectionData struct {
	Message string
}

func ProfilePasswordSection(d *ProfilePasswordSectionData) e.Node {
	return e.Form(
		e.HXPost("/r/change-password"),
		e.Div(
			e.Class("flex flex-col gap-y-4"),
			e.H2(e.Class("text-lg font-medium"), e.Raw("Password")),
			e.H3(e.Class("font-medium"), e.Raw("Change your password")),
			e.Div(
				e.Class("flex flex-col gap-y-2 max-w-md"),
				e.Label(e.Class("text-sm font-medium"), e.Raw("Current password")),
				e.Input(e.Type("password"), e.Name("currentPassword"), e.Class(cn(cnInput, "")), e.AutoComplete("current-password"), e.Required()),
				e.Label(e.Class("text-sm font-medium"), e.Raw("New password")),
				e.Input(e.Type("password"), e.Name("newPassword"), e.Class(cn(cnInput, "")), e.AutoComplete("new-password"), e.Required()),
				e.If(d.Message != "", e.Div(e.Class("text-red-400"), e.Text(d.Message))),
				e.Button(e.Type("submit"), e.Class(cn(btn, btnPrimary)), e.Raw("Save changes")),
			),
		),
	)
}
