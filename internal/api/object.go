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
	ETag        string    `json:"etag"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
}

func newObjectResponse(o *object.Object) ObjectResponse {
	return ObjectResponse{
		Key:         o.Key,
		ContentType: o.ContentType,
		Size:        o.Size,
		ETag:        o.ETag,
		CreatedAt:   o.CreatedAt,
	}
}

func handleObjectGet(c jug.Context) {
	query := c.Request().URL.Query()
	if query.Has(queryArchiveId) {
		handleGetArchive(c)
	} else {
		handleGetObject(c)
	}
}

func handleGetObject(c jug.Context) {
	o, ok := objectFilter(c)
	if !ok {
		return
	}
	c.Status(200)
	c.SetHeader("Content-Length", strconv.FormatInt(int64(o.Size), 10))
	c.SetHeader("Content-Type", o.ContentType)

	if err := object.Write(c, o, c.Writer()); err != nil {
		log.Printf("unable to write object: %v", err)
	}
}

func handleObjectPost(c jug.Context) {
	query := c.Request().URL.Query()
	if query.Has(queryArchives) {
		handleCreateArchive(c)
	} else if query.Get(queryArchiveId) != "" {
		handleCompleteArchive(c)
	} else if query.Has(queryUploads) {
		handleCreateMultipartUpload(c)
	} else if query.Get(queryUploadId) != "" {
		handleCompleteMultipartUpload(c)
	} else {
		c.Status(405)
	}
}

func handleObjectPut(c jug.Context) {
	if c.Query(queryArchiveId) != "" {
		handleAddArchiveEntries(c)
	} else if c.Query(queryUploadId) != "" {
		handleUploadPart(c)
	} else {
		handleCreateObject(c)
	}
}

func handleObjectDelete(c jug.Context) {
	if c.Query(queryArchiveId) != "" {
		handleAbortArchive(c)
	} else if c.Query(queryUploadId) != "" {
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

	o, err := uc.CreateObject(c, b, object.CreateCommand{
		Key:         key,
		ContentType: contentType,
		Data:        d,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(204)
	c.SetHeader("ETag", o.ETag)
}

func handleDeleteObject(c jug.Context) {
	o, ok := objectFilter(c)
	if !ok {
		return
	}
	b := contextGetBucket(c)

	if err := uc.DeleteObject(c, b, o); err != nil {
		handleError(c, err)
		return
	}

	c.RespondNoContent()
}
