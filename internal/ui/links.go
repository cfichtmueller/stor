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

type bucketLinks struct {
	Objects    string
	Properties string
	Settings   string
}

func newBucketLinks(bucketName string) bucketLinks {
	base := bucketsLink + "/" + bucketName
	return bucketLinks{
		Objects:    base + "/objects",
		Properties: base + "/properties",
		Settings:   base + "/settings",
	}
}
