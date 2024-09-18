// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/disk"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/domain/user"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleDashboardPage(c jug.Context) {
	info, err := disk.GetInfo(config.DataDir)
	if err != nil {
		c.HandleError(err)
		return
	}
	bucketStats, err := bucket.GetStats(c)
	if err != nil {
		c.HandleError(err)
		return
	}
	chunkStats, err := chunk.GetStats(c)
	if err != nil {
		c.HandleError(err)
		return
	}

	must("render dashboard page", c, ui.RenderDashboardPage(c.Writer(), ui.DashboardData{
		DiskInfo:    info,
		StorageSize: chunkStats.TotalSize,
		BucketStats: bucketStats,
	}))
}

func handleBucketsPage(c jug.Context) {
	b, err := bucket.FindMany(c, &bucket.Filter{})
	if !must("find buckets", c, err) {
		return
	}
	must("render buckets page", c, ui.RenderBucketsPage(c.Writer(), b))
}

func handleBucketPage(c jug.Context) {
	b := contextGetBucket(c)
	c.Status(302)
	c.SetHeader("Location", "/u/buckets/"+b.Name+"/files")
}

func handleBucketObjectsPage(c jug.Context) {
	b := contextGetBucket(c)
	objects, err := object.List(c, b.Name)
	if !must("find objects", c, err) {
		return
	}
	must("render bucket objects page", c, ui.RenderBucketObjectsPage(c.Writer(), b, objects))
}

func handleBucketSettingsPage(c jug.Context) {
	b := contextGetBucket(c)
	must("render bucket settings page", c, ui.RenderBucketSettingsPage(c.Writer(), b))
}

func handleAdminPage(c jug.Context) {
	hxRedirect(c, "/u/admin/users")
}

func handleUsersPage(c jug.Context) {
	u, err := user.List(c)
	if err != nil {
		c.HandleError(err)
		return
	}
	must("render users page", c, ui.RenderUsersPage(c.Writer(), ui.UsersPageData{
		Users: u,
	}))
}

func handleApiKeysPage(c jug.Context) {
	keys, err := apikey.List(c)
	if err != nil {
		c.HandleError(err)
		return
	}

	must("render api keys page", c, ui.RenderApiKeysPage(c.Writer(), ui.ApiKeysPageData{
		Keys: keys,
	}))
}
