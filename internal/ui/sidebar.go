package ui

import "github.com/cfichtmueller/goparts/e"

type SidebarItem struct {
	Active bool
	Icon   e.Node
	Title  string
	Link   string
}

func Sidebar(items ...SidebarItem) e.Node {
	return e.Aside(
		e.Class("p-10 pt-0 ps-6 min-w-48"),
		e.Div(
			e.Class("flex items-center p-4 gap-4"),
			e.Img(e.Attr("src", "/img/icon.png?v=1726228022"), e.Class("h-4")),
			e.Span(e.Class("text-sm font-medium"), e.Raw("STOR")),
		),
		e.Nav(
			e.Class("flex flex-col gap-y-1"),
			e.Mapf(items, func(i SidebarItem) e.Node {
				return e.Div(
					e.HXBoost(),
					e.If(i.Active, e.Div(
						e.Class("inline-flex items-center text-sm font-medium py-0.5 whitespace-nowrap rounded-md h-9 w-full px-4 hover:bg-neutral-200 bg-neutral-200 cursor-pointer"),
						i.Icon,
						e.Raw(i.Title),
					)),
					e.If(!i.Active, e.A(
						e.Class("inline-flex items-center text-sm font-medium py-0.5 whitespace-nowrap rounded-md h-9 w-full px-4 hover:bg-neutral-200"),
						e.Href(i.Link),
						i.Icon,
						e.Raw(i.Title),
					)),
				)
			}),
		),
	)
}
