// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/config"
)

func Configure() *srv.Server {

	server := srv.NewServer().SetTrustRemoteIdHeaders(config.TrustProxies).Use(srv.LoggingMiddleware())

	server.GET("", handleListBuckets, authenticatedFilter)

	bucketGroup := server.Group("/{bucketName}", authenticatedFilter)
	bucketGroup.POST("", handleBucketPost, bucketFilter)
	bucketGroup.PUT("", handleCreateBucket)
	bucketGroup.GET("", handleListObjects, bucketFilter)
	bucketGroup.DELETE("", handleDeleteBucket, bucketFilter)

	objectGroup := server.Group("/{bucketName}/{objectKey...}")
	objectGroup.GET("", handleObjectGet)
	objectGroup.POST("", handleObjectPost, authenticatedFilter, bucketFilter)
	objectGroup.PUT("", handleObjectPut, authenticatedFilter, bucketFilter)
	objectGroup.DELETE("", handleObjectDelete, authenticatedFilter, bucketFilter)

	return server
}
