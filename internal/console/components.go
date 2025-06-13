// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/disk"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/ui"
)

//
// API Key
//

func handleRenderApiKeySheet(c *srv.Context) (e.Node, error) {
	key := contextGetApiKey(c)
	return ui.ApiKeySheet(key), nil
}

func handleRenderApiKeysTable(c *srv.Context) (e.Node, error) {
	keys, err := apikey.List(c)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return ui.ApiKeysEmptyState(), nil
	}
	return ui.ApiKeysTable(keys), nil
}

func handleRenderDeleteApiKeyDialog(c *srv.Context) (e.Node, error) {
	key := contextGetApiKey(c)
	return ui.DeleteApiKeyDialog(key), nil
}

//
// Bucket
//

func handleRenderBucketsTable(c *srv.Context) (e.Node, error) {
	b, err := bucket.FindMany(c, &bucket.Filter{})
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return ui.BucketEmptyState(), nil
	}
	return ui.BucketsTable(b), nil
}

//
// Dashboard
//

func handleRenderDashboardMetrics(c *srv.Context) (e.Node, error) {
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

	return ui.DashboardMetrics(ui.DashboardData{
		DiskInfo:    info,
		StorageSize: chunkStats.TotalSize,
		BucketStats: bucketStats,
	}), nil
}

func renderNodeFn(f func() e.Node) func(c *srv.Context) *srv.Response {
	return renderNode(func(c *srv.Context) (e.Node, error) {
		return f(), nil
	})
}
