// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func BucketsPage(buckets []*bucket.Bucket) e.Node {
	empty := len(buckets) == 0
	return ListPageLayout(
		"Buckets",
		appSidebar(app_sidebar_active_buckets),
		e.If(empty, BucketsEmptyState()),
		e.If(!empty, BucketsTable(buckets)),
	)
}

func BucketsEmptyState() e.Node {
	return e.Div(
		e.Class("flex flex-col justify-center items-center min-h-96"),
		e.Img(
			e.Class("h-32 mb-8"),
			e.AriaHidden(),
			e.Attr("src", "/img/empty.png"),
		),
		e.Button(
			e.Class(cn(btn, btnPrimary)),
			e.HXGet("/c/create-bucket-dialog"),
			e.HXTarget("body"),
			e.HXSwap("beforeend"),
			e.Raw("Create your first bucket"),
		),
	)
}

func BucketsTable(buckets []*bucket.Bucket) e.Node {
	return Table(
		TableHeader(
			TableHead("", e.Raw("Name")),
			TableHead("text-right", e.Raw("Objects")),
			TableHead("text-right", e.Raw("Size")),
			TableHead("flex justify-end", e.Button(
				e.Class(cn(btn, btnPrimary)),
				e.HXGet("/c/create-bucket-dialog"),
				e.HXTarget("body"),
				e.HXSwap("beforeend"),
				IconPlus,
				e.Span(
					srOnly(),
					e.Raw("Create"),
				),
			)),
		),
		TableBody(
			e.Mapf(buckets, func(b *bucket.Bucket) e.Node {
				return TableRow(
					e.Td(
						e.Class("p-2 w-1 align-middle whitespace-nowrap"),
						e.HXBoost(),
						e.A(
							e.Attr("href", fmt.Sprintf("/u/buckets/%s/objects", b.Name)),
							e.Raw(b.Name),
						),
					),
					e.Td(
						e.Class("text-right"),
						e.Raw(formatInt(int(b.Objects))),
					),
					e.Td(
						e.Class("text-right"),
						e.Raw(formatBytes(b.Size)),
					),
				)
			}),
		),
	)
}
