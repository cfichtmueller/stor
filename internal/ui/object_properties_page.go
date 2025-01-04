// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/object"
)

func ObjectPropertiesPage(b *bucket.Bucket, o *object.Object) e.Node {
	links := NewBucketLinks(b.Name)
	return LoggedInLayout(
		appSidebar(app_sidebar_active_buckets),
		PathBreadcrumbs(links, b, o.Key),
		PageTitle(""),
		e.Div(
			e.Class("flex flex-col w-full border rounded-md bg-white"),
			NavTabs(
				ObjectsNavTab(links.Folder(object.PathPrefix(o.Key, "/")), false),
				PropertiesNavTab("", true),
			),
			e.Div(
				e.Class("p-2"),
				Details("",
					Detail("Key", o.Key),
					Detail("Size", formatBytes(o.Size)),
					Detail("Created at", formatDateTime(o.CreatedAt)),
				),
				e.Div(
					e.Class("flex justify-end gap-x-2"),
					e.Button(
						e.Class(cn(btn, "shadow")),
						e.A(
							e.Href(OpenObjectLink(b.Name, o.Key)),
							e.TargetBlank(),
							e.Raw("Open"),
						),
					),
					e.Button(
						e.Class(cn(btn, "shadow")),
						e.A(
							e.Href(DownloadObjectLink(b.Name, o.Key)),
							e.TargetBlank(),
							e.Raw("Download"),
						),
					),
				),
			),
		),
	)
}
