// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
)

func Configure() jug.Engine {

	engine := jug.New()

	bucketGroup := engine.Group("/:bucketName", authenticatedFilter, bucketFilter)
	bucketGroup.GET("", handleListObjects)
	objectGroup := bucketGroup.Group("/*objectKey")
	objectGroup.GET("", objectFilter, handleGetObject)
	objectGroup.PUT("", handleCreateObject)
	objectGroup.DELETE("", objectFilter, handleDeleteObject)

	return engine
}
