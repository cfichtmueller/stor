// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/disk"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

type DashboardData struct {
	DiskInfo    disk.Info
	StorageSize uint64
	BucketStats bucket.Stats
}

func DashboardPage(d DashboardData) e.Node {
	return LoggedInLayout(
		appSidebar(app_sidebar_active_dashboard),
		e.Div(
			e.Class("w-full flex-grow px-2"),
			PageHeader("Dashboard", ""),
			e.Div(
				e.Class("rounded-lg bg-white border flex flex-col gap-y-4 p-4 w-full"),
				e.HXTrigger("bucketsUpdated from:body"),
				e.HXGet("/c/dashboard-metrics"),
				DashboardMetrics(d),
			),
		),
	)
}

func DashboardMetrics(d DashboardData) e.Node {
	return e.Group(
		e.Div(
			e.Class("grid gap-4 md:grid-cols-2 lg:grid-cols-4"),
			MetricCard("Total space", nil, formatBytes(int64(d.DiskInfo.Total)), ""),
			MetricCard("Used Space", nil, formatBytes(int64(d.DiskInfo.Used)), ""),
			MetricCard("Free Space", nil, formatBytes(int64(d.DiskInfo.Free)), ""),
			MetricCard("Storage Space", nil, formatBytes(int64(d.StorageSize)), ""),
		),
		e.Div(
			e.Class("grid gap-4 md:grid-cols-2"),
			MetricCard("Buckets", nil, formatInt(d.BucketStats.Count), ""),
			MetricCard("Objects", nil, formatInt(d.BucketStats.TotalObjects), ""),
		),
		e.If(d.BucketStats.Count == 0, e.Div(
			e.Class("w-full min-h96 flex justify-center items-center"),
			e.Button(
				e.Class(cn(btn, btnPrimary)),
				e.HXGet("/c/create-bucket-dialog"),
				e.HXTarget("body"),
				e.HXSwap("beforeend"),
				e.Raw("Create your first bucket"),
			),
		)),
	)
}
