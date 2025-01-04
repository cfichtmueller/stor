// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

const (
	bucket_navtabs_active_objects    = "objects"
	bucket_navtabs_active_properties = "properties"
	bucket_navtabs_active_settings   = "settings"
)

func BucketNavTabs(links *BucketLinks, active string) e.Node {
	return NavTabs(
		ObjectsNavTab(links.Objects, active == bucket_navtabs_active_objects),
		&NavLink{
			Title:  "Properties",
			Link:   links.Properties,
			Icon:   IconSlidersHorizontal,
			Active: active == bucket_navtabs_active_properties,
		},
		&NavLink{
			Title:  "Settings",
			Link:   links.Settings,
			Icon:   IconCog,
			Active: active == bucket_navtabs_active_settings,
		},
	)
}
