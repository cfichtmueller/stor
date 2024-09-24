// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/disk"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

type DashboardData struct {
	DiskInfo    disk.Info
	StorageSize uint64
	BucketStats bucket.Stats
}

type dashboardMetricsModel struct {
	TotalSpaceCard   MetricCardModel
	FreeSpaceCard    MetricCardModel
	UsedSpaceCard    MetricCardModel
	StorageSpaceCard MetricCardModel
	BucketCountCard  MetricCardModel
	ObjectCountCard  MetricCardModel
	NoBuckets        bool
}

func newDashboardMetricsModel(d DashboardData) dashboardMetricsModel {
	return dashboardMetricsModel{
		TotalSpaceCard: MetricCardModel{
			Title: "Total Space",
			Value: formatBytes(int64(d.DiskInfo.Total)),
		},
		UsedSpaceCard: MetricCardModel{
			Title: "Used Space",
			Value: formatBytes(int64(d.DiskInfo.Used)),
		},
		FreeSpaceCard: MetricCardModel{
			Title: "Free Space",
			Value: formatBytes(int64(d.DiskInfo.Free)),
		},
		StorageSpaceCard: MetricCardModel{
			Title: "Storage Space",
			Value: formatBytes(int64(d.StorageSize)),
		},
		BucketCountCard: MetricCardModel{
			Title: "Buckets",
			Value: formatInt(d.BucketStats.Count),
		},
		ObjectCountCard: MetricCardModel{
			Title: "Objects",
			Value: formatInt(d.BucketStats.TotalObjects),
		},
		NoBuckets: d.BucketStats.Count == 0,
	}
}

func RenderDashboardMetrics(w io.Writer, d DashboardData) error {
	return renderTemplate(w, "DashboardMetrics", newDashboardMetricsModel(d))
}
