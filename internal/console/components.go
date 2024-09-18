// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"io"

	"github.com/cfichtmueller/jug"
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

func handleRenderApiKeySheet(c jug.Context) {
	key := contextGetApiKey(c)
	must("render api key sheet", c, ui.RenderApiKeySheet(c.Writer(), key))
}

func handleRenderApiKeysTable(c jug.Context) {
	keys, err := apikey.List(c)
	if !must("find api keys", c, err) {
		return
	}

	if len(keys) == 0 {
		must("render api keys empty state", c, ui.RenderApiKeysEmptyState(c.Writer()))
		return
	}

	must("render api keys table", c, ui.RenderApiKeysTable(c.Writer(), keys))
}

func handleRenderDeleteApiKeyDialog(c jug.Context) {
	key := contextGetApiKey(c)
	must("render delete api key dialog", c, ui.RenderDeleteApiKeyDialog(c.Writer(), key))
}

//
// Bucket
//

func handleRenderBucketsTable(c jug.Context) {
	b, err := bucket.FindMany(c, &bucket.Filter{})
	if !must("find buckets", c, err) {
		return
	}

	if len(b) == 0 {
		must("render buckets table", c, ui.RenderBucketsEmptyState(c.Writer()))
		return
	}

	must("render buckets table", c, ui.RenderBucketsTable(c.Writer(), b))
}

//
// Dashboard
//

func handleRenderDashboardMetrics(c jug.Context) {
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

	must("render dashboard metrics", c, ui.RenderDashboardMetrics(c.Writer(), ui.DashboardData{
		DiskInfo:    info,
		StorageSize: chunkStats.TotalSize,
		BucketStats: bucketStats,
	}))
}

func uiRenderFn(name string, f func(w io.Writer) error) func(c jug.Context) {
	return func(c jug.Context) {
		must("render "+name, c, f(c.Writer()))
	}
}
