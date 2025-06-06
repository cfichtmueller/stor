// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log/slog"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
)

type CreateMultipartUploadResult struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	UploadId string `json:"uploadId"`
}

func handleCreateMultipartUpload(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	contentType := c.Request().Header.Get("Content-Type")

	slog.Info("create multipart upload", "bucket", b.Name, "key", key, "content-type", contentType)

	c.RespondOk(CreateMultipartUploadResult{
		Bucket:   b.Name,
		Key:      key,
		UploadId: domain.RandomId(),
	})
}

func handleUploadPart(c jug.Context) {
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		handleError(c, ec.InvalidArgument)
	}
	partNumber, err := c.IntQuery("part-number")
	if err != nil {
		handleError(c, ec.InvalidArgument)
		return
	}

	slog.Info("upload part", "upload", uploadId, "part", partNumber)

	c.Status(200)
	c.SetHeader("ETag", domain.RandomId())
}

type PartReference struct {
	ETag       string `json:"etag"`
	PartNumber int    `json:"partNumber"`
}

type CompleteMultipartUploadRequest struct {
	Parts []PartReference `json:"parts"`
}

type CompleteMultipartUploadResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	ETag   string `json:"etag"`
}

func handleCompleteMultipartUpload(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		handleError(c, ec.InvalidArgument)
		return
	}
	var req CompleteMultipartUploadRequest
	if !c.MustBindJSON(&req) {
		return
	}

	slog.Info("complete multipart upload", "bucket", b.Name, "key", key, "upload", uploadId)

	c.RespondOk(CompleteMultipartUploadResult{
		Bucket: b.Name,
		Key:    key,
		ETag:   domain.RandomId(),
	})
}

func handleAbortMultipartUpload(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		handleError(c, ec.InvalidArgument)
		return
	}

	slog.Info("abort multipart upload", "bucket", b.Name, "key", key, "upload", uploadId)

	c.RespondNoContent()
}
