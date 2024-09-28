// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

const (
	dashboardLink = "/u"
	bucketsLink   = "/u/buckets"
	adminLink     = "/u/admin"
	usersLink     = "/u/admin/users"
	apiKeysLink   = "/u/admin/api-keys"
)

type BucketLinks struct {
	Objects    string
	Properties string
	Settings   string
}

func (l *BucketLinks) Folder(prefix string) string {
	return l.Objects + "?prefix=" + prefix
}

func NewBucketLinks(bucketName string) *BucketLinks {
	base := bucketsLink + "/" + bucketName
	return &BucketLinks{
		Objects:    base + "/objects",
		Properties: base + "/properties",
		Settings:   base + "/settings",
	}
}
