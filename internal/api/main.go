// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/ec"
)

func Configure() jug.Engine {

	engine := jug.Default()

	engine.GET("", authenticatedFilter, handleListBuckets)
	bucketGroup := engine.Group("/:bucketName", authenticatedFilter)
	bucketGroup.POST("", bucketFilter, handleBucketPost)
	bucketGroup.PUT("", handleCreateBucket)
	bucketGroup.GET("", bucketFilter, handleListObjects)
	bucketGroup.DELETE("", bucketFilter, handleDeleteBucket)
	objectGroup := engine.Group("/:bucketName/*objectKey")
	objectGroup.GET("", handleObjectGet)
	objectGroup.POST("", authenticatedFilter, bucketFilter, handleObjectPost)
	objectGroup.PUT("", authenticatedFilter, bucketFilter, handleObjectPut)
	objectGroup.DELETE("", authenticatedFilter, bucketFilter, handleObjectDelete)

	return engine
}

func handleError(ctx jug.Context, err error) {
	e, ok := err.(*ec.Error)
	if !ok {
		ctx.HandleError(err)
		ctx.Abort()
		return
	}

	b, err := json.Marshal(e)
	if err != nil {
		ctx.HandleError(err)
		ctx.Abort()
		return
	}

	ctx.Data(e.StatusCode, "application/json", b)
}
