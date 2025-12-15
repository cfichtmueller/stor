// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"slices"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

type ProfilePageData struct {
	User                 *user.User
	AuthenticatedSession *session.Session
	Sessions             []*session.Session
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
				ProfileSessionsSection(&ProfileSessionsSectionData{
					AuthenticatedSession: d.AuthenticatedSession,
					Sessions:             d.Sessions,
				}),
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
		ProfileSection(
			"Password",
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

type ProfileSessionsSectionData struct {
	AuthenticatedSession *session.Session
	Sessions             []*session.Session
}

func ProfileSessionsSection(d *ProfileSessionsSectionData) e.Node {
	sessions := d.Sessions[:]
	slices.SortFunc(sessions, func(a, b *session.Session) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})
	return ProfileSection(
		"Sessions",
		e.H3(e.Class("font-medium"), e.Raw("Active sessions")),
		e.Div(
			e.Class("flex flex-col gap-y-2 max-w-md"),
			e.Ul(
				e.Mapf(sessions, func(s *session.Session) e.Node {
					return e.Li(
						e.Class("grid grid grid-cols-9 gap-x-2 items-center py-2"),
						e.Id("session-"+s.ID),
						e.Div(e.Class("col-span-3 text-sm text-neutral-500 whitespace-nowrap"), e.Raw(s.IpAddress)),
						e.Div(e.Class("col-span-4 text-sm text-neutral-500 whitespace-nowrap"), e.Raw(formatDateTime(s.LastSeenAt))),
						e.Button(
							e.Type("button"),
							e.Class(cn(btn, btnPrimary, "col-span-2")),
							e.HXPost("/r/logout-session?session="+s.ID),
							e.HXTarget("#session-"+s.ID),
							e.Raw("Logout"),
						),
						e.If(
							s.ID == d.AuthenticatedSession.ID,
							e.Div(e.Class("grid-col-span-9 text-sm  whitespace-nowrap"), e.Raw("This is your current session")),
						),
					)
				}),
			),
		),
	)
}

func ProfileSection(title string, children ...e.Node) e.Node {
	return e.Div(
		e.Class("flex flex-col gap-y-4 py-4"),
		e.H2(e.Class("text-lg font-medium"), e.Raw(title)),
		e.Div(children...),
	)
}
