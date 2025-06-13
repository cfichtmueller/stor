// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/disk"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleDashboardPage(c *srv.Context) *srv.Response {
	info, err := disk.GetInfo(config.DataDir)
	if err != nil {
		return responseFromError(err)
	}
	bucketStats, err := bucket.GetStats(c)
	if err != nil {
		return responseFromError(err)
	}
	chunkStats, err := chunk.GetStats(c)
	if err != nil {
		return responseFromError(err)
	}

	return nodeResponseWithShell(c, ui.DashboardPage(ui.DashboardData{
		DiskInfo:    info,
		StorageSize: chunkStats.TotalSize,
		BucketStats: bucketStats,
	}))
}

func handleBucketsPage(c *srv.Context) *srv.Response {
	b, err := bucket.FindMany(c, &bucket.Filter{})
	if err != nil {
		return responseFromError(err)
	}

	return nodeResponseWithShell(c, ui.BucketsPage(b))
}

func handleBucketPage(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	return hxRedirect(c, "/u/buckets/"+b.Name+"/objects")
}

func handleBucketObjectsPage(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	delimiter := "/"
	prefix := c.Query("prefix")
	prefixLen := len(prefix)
	r, err := uc.ObjectPrefixSearch(c, b, delimiter, prefix, "", 1000)
	if err != nil {
		return responseFromError(err)
	}
	bucketLinks := ui.NewBucketLinks(b.Name)
	objects := make([]ui.ObjectData, 0, len(r.CommonPrefixes)+len(r.Objects)+1)
	if pathParts := object.SplitPath(prefix, delimiter); len(pathParts) > 0 {
		objects = append(objects, ui.ObjectData{
			Key:  "..",
			Href: bucketLinks.Folder(object.JoinPath(pathParts[:len(pathParts)-1], delimiter)),
		})
	}
	for _, p := range r.CommonPrefixes {
		objects = append(objects, ui.ObjectData{
			Key:  p[prefixLen:],
			Href: bucketLinks.Folder(p),
		})
	}
	for _, o := range r.Objects {
		objects = append(objects, ui.ObjectData{
			Key:  o.Key[prefixLen:],
			Size: o.Size,
			Href: bucketLinks.Object(o.Key),
		})
	}
	return nodeResponseWithShell(c, ui.BucketObjectsPage(ui.BucketObjectsPageData{
		Bucket:  b,
		Prefix:  prefix,
		Objects: objects,
	}))
}

func handleBucketPropertiesPage(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	return nodeResponseWithShell(c, ui.BucketPropertiesPage(b))
}

func handleBucketSettingsPage(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	return nodeResponseWithShell(c, ui.BucketSettingsPage(b))
}

func handleObjectPage(c *srv.Context) *srv.Response {
	key := c.Query("key")
	b := contextGetBucket(c)
	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil && errors.Is(err, ec.NoSuchKey) {
		return nodeResponseWithShell(c, ui.NotFoundPage())
	}
	if err != nil {
		return responseFromError(err)
	}
	return nodeResponseWithShell(c, ui.ObjectPropertiesPage(b, o))
}

func handleUsersPage(c *srv.Context) *srv.Response {
	u, err := user.List(c)
	if err != nil {
		return responseFromError(err)
	}
	return nodeResponseWithShell(c, ui.UsersPage(&ui.UsersPageData{
		Users: u,
	}))
}

func handleApiKeysPage(c *srv.Context) *srv.Response {
	keys, err := apikey.List(c)
	if err != nil {
		return responseFromError(err)
	}

	return nodeResponseWithShell(c, ui.ApiKeysPage(&ui.ApiKeysPageData{
		Keys: keys,
	}))
}
