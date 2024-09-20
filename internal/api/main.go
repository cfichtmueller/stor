// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
)

func Configure() jug.Engine {

	engine := jug.Default()

	bucketGroup := engine.Group("/:bucketName", authenticatedFilter)
	bucketGroup.POST("", bucketFilter, handleBucketPost)
	bucketGroup.PUT("", handleCreateBucket)
	bucketGroup.GET("", bucketFilter, handleListObjects)
	objectGroup := bucketGroup.Group("/*objectKey", bucketFilter)
	objectGroup.GET("", objectFilter, handleGetObject)
	objectGroup.PUT("", handleCreateObject)
	objectGroup.DELETE("", objectFilter, handleDeleteObject)

	return engine
}
