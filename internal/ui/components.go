// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

type BreadcrumbsModel struct {
	Crumbs []*BreadcrumbModel
}

func (m *BreadcrumbsModel) Add(c *BreadcrumbModel) {
	m.Crumbs = append(m.Crumbs, c)
}

func (m *BreadcrumbsModel) Last() *BreadcrumbModel {
	if len(m.Crumbs) == 0 {
		return nil
	}
	return m.Crumbs[len(m.Crumbs)-1]
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
