// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

type ObjectData struct {
	Key  string
	Size int64
	Href string
}

type objectModel struct {
	Key  string
	Size string
	Href string
}

func newObjectModel(d ObjectData) objectModel {
	size := ""
	if d.Size > 0 {
		size = formatBytes(d.Size)
	}
	return objectModel{
		Key:  d.Key,
		Size: size,
		Href: d.Href,
	}
}

func newObjectNavTabs(objectsLink, active string) *NavTabsModel {
	return &NavTabsModel{
		Tabs: []*NavLink{
			newObjectsTab(objectsLink, active == "objects"),
			newPropertiesTab("", active == "properties"),
		},
	}
}
