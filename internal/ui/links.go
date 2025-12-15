// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"net/url"
)

const (
	dashboardLink = "/u"
	bucketsLink   = "/u/buckets"
	adminLink     = "/u/admin"
	usersLink     = "/u/admin/users"
	apiKeysLink   = "/u/admin/api-keys"
	profileLink   = "/u/profile"
)

type BucketLinks struct {
	base       string
	Objects    string
	Properties string
	Settings   string
}

func NewBucketLinks(bucketName string) *BucketLinks {
	base := bucketsLink + "/" + bucketName
	return &BucketLinks{
		base:       base,
		Objects:    base + "/objects",
		Properties: base + "/properties",
		Settings:   base + "/settings",
	}
}

func (l *BucketLinks) Folder(prefix string) string {
	return l.Objects + "?prefix=" + prefix
}

func (l *BucketLinks) Object(key string) string {
	return l.base + "/object?key=" + key
}

func DownloadObjectLink(bucket, key string) string {
	return fmt.Sprintf("/download?bucket=%s&key=%s", bucket, url.QueryEscape(key))
}

func OpenObjectLink(bucket, key string) string {
	return fmt.Sprintf("/open?bucket=%s&key=%s", bucket, url.QueryEscape(key))
}
