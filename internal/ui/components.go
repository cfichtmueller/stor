// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/stor/internal/domain/object"

type BreadcrumbsModel struct {
	Crumbs []*BreadcrumbModel
}

func NewBreadcrumbs() *BreadcrumbsModel {
	return &BreadcrumbsModel{
		Crumbs: make([]*BreadcrumbModel, 0),
	}
}

func (m *BreadcrumbsModel) AddLink(title, link string) *BreadcrumbsModel {
	return m.add(&BreadcrumbModel{
		Title: title,
		Link:  link,
	})
}

func (m *BreadcrumbsModel) AddTitle(title string) *BreadcrumbsModel {
	return m.add(&BreadcrumbModel{Title: title})
}

func (m *BreadcrumbsModel) Last() *BreadcrumbModel {
	if len(m.Crumbs) == 0 {
		return nil
	}
	return m.Crumbs[len(m.Crumbs)-1]
}

func (m *BreadcrumbsModel) add(c *BreadcrumbModel) *BreadcrumbsModel {
	if len(m.Crumbs) > 0 {
		m.Crumbs = append(m.Crumbs, &BreadcrumbModel{Separator: true})
	}
	m.Crumbs = append(m.Crumbs, c)
	return m
}

type BreadcrumbModel struct {
	Separator bool
	Title     string
	Link      string
}

type DetailsModel struct {
	Title   string
	Details []DetailModel
}

type DetailModel struct {
	Title string
	Value string
}

func newBucketBreadcrumbs(name string) *BreadcrumbsModel {
	return NewBreadcrumbs().AddLink("Buckets", bucketsLink).AddTitle(name)
}

func addPathCrumbs(b *BreadcrumbsModel, links *BucketLinks, key string) {
	prefix := ""
	for _, f := range object.SplitPath(key, "/") {
		prefix = prefix + f + "/"
		b.AddLink(f, links.Folder(prefix))
	}
}

func newObjectsTab(link string, active bool) *NavLink {
	return &NavLink{
		Link:   link,
		Active: active,
		Title:  "Objects",
		Icon:   "files",
	}
}

func newPropertiesTab(link string, active bool) *NavLink {
	return &NavLink{
		Link:   link,
		Active: active,
		Title:  "Properties",
		Icon:   "sliders-horizontal",
	}
}

func newSettingsTab(link string, active bool) *NavLink {
	return &NavLink{
		Link:   link,
		Active: active,
		Title:  "Settings",
		Icon:   "cog",
	}
}
