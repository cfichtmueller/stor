// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/jug"
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

func handleDashboardPage(c jug.Context) (e.Node, error) {
	info, err := disk.GetInfo(config.DataDir)
	if err != nil {
		return nil, err
	}
	bucketStats, err := bucket.GetStats(c)
	if err != nil {
		return nil, err
	}
	chunkStats, err := chunk.GetStats(c)
	if err != nil {
		return nil, err
	}

	return ui.DashboardPage(ui.DashboardData{
		DiskInfo:    info,
		StorageSize: chunkStats.TotalSize,
		BucketStats: bucketStats,
	}), nil
}

func handleBucketsPage(c jug.Context) (e.Node, error) {
	b, err := bucket.FindMany(c, &bucket.Filter{})
	if err != nil {
		return nil, err
	}

	return ui.BucketsPage(b), nil
}

func handleBucketPage(c jug.Context) {
	b := contextGetBucket(c)
	c.Status(302)
	c.SetHeader("Location", "/u/buckets/"+b.Name+"/files")
}

func handleBucketObjectsPage(c jug.Context) (e.Node, error) {
	b := contextGetBucket(c)
	delimiter := "/"
	prefix := c.Query("prefix")
	prefixLen := len(prefix)
	r, err := uc.ObjectPrefixSearch(c, b, delimiter, prefix, "", 1000)
	if err != nil {
		return nil, err
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
	return ui.BucketObjectsPage(ui.BucketObjectsPageData{
		Bucket:  b,
		Prefix:  prefix,
		Objects: objects,
	}), nil
}

func handleBucketPropertiesPage(c jug.Context) (e.Node, error) {
	b := contextGetBucket(c)
	return ui.BucketPropertiesPage(b), nil
}

func handleBucketSettingsPage(c jug.Context) (e.Node, error) {
	b := contextGetBucket(c)
	return ui.BucketSettingsPage(b), nil
}

func handleObjectPage(c jug.Context) (e.Node, error) {
	key := c.Query("key")
	b := contextGetBucket(c)
	o, err := object.FindOne(c, b.Name, key, false)
	if err != nil && errors.Is(err, ec.NoSuchKey) {
		return ui.NotFoundPage(), nil
	}
	if err != nil {
		return nil, err
	}
	return ui.ObjectPropertiesPage(b, o), nil
}

func handleUsersPage(c jug.Context) (e.Node, error) {
	u, err := user.List(c)
	if err != nil {
		return nil, err
	}
	return ui.UsersPage(&ui.UsersPageData{
		Users: u,
	}), nil
}

func handleApiKeysPage(c jug.Context) (e.Node, error) {
	keys, err := apikey.List(c)
	if err != nil {
		return nil, err
	}

	return ui.ApiKeysPage(&ui.ApiKeysPageData{
		Keys: keys,
	}), nil
}
