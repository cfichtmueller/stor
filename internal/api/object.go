// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log"
	"strconv"
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/uc"
)

type ObjectResponse struct {
	Key         string    `json:"key"`
	ContentType string    `json:"contentType"`
	Size        uint64    `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
}

func newObjectResponse(o *object.Object) ObjectResponse {
	return ObjectResponse{
		Key:         o.Key,
		ContentType: o.ContentType,
		Size:        o.Size,
		CreatedAt:   o.CreatedAt,
	}
}

func handleGetObject(c jug.Context) {
	o := contextGetObject(c)
	c.Status(200)
	c.SetHeader("Content-Length", strconv.FormatInt(int64(o.Size), 10))
	c.SetHeader("Content-Type", o.ContentType)

	if err := object.Write(c, o, c.Writer()); err != nil {
		log.Printf("unable to write object: %v", err)
	}
}

func handleObjectPost(c jug.Context) {
	if c.Request().URL.Query().Has(queryUploads) {
		handleCreateMultipartUpload(c)
	} else if c.Query(queryUploadId) != "" {
		handleCompleteMultipartUpload(c)
	} else {
		c.Status(405)
	}
}

func handleObjectPut(c jug.Context) {
	if c.Query(queryUploadId) != "" {
		handleUploadPart(c)
	} else {
		handleCreateObject(c)
	}
}

func handleObjectDelete(c jug.Context) {
	if c.Query(queryUploadId) != "" {
		handleAbortMultipartUpload(c)
	} else {
		handleDeleteObject(c)
	}
}

func handleCreateObject(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	d, err := c.GetRawData()
	if err != nil {
		handleError(c, err)
		return
	}

	if err := uc.CreateObject(c, b, object.CreateCommand{
		Key:         key,
		ContentType: contentType,
		Data:        d,
	}); err != nil {
		handleError(c, err)
		return
	}

	c.RespondNoContent()
}

func handleDeleteObject(c jug.Context) {
	b := contextGetBucket(c)
	o := contextGetObject(c)

	if err := uc.DeleteObject(c, b, o); err != nil {
		handleError(c, err)
		return
	}

	c.RespondNoContent()
}
