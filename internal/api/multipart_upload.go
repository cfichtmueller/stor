// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log/slog"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
)

type CreateMultipartUploadResult struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	UploadId string `json:"uploadId"`
}

func handleCreateMultipartUpload(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	key, r := contextGetObjectKey(c)
	if r != nil {
		return r
	}
	contentType := c.Request().Header.Get("Content-Type")

	slog.Info("create multipart upload", "bucket", b.Name, "key", key, "content-type", contentType)

	return srv.Respond().Json(CreateMultipartUploadResult{
		Bucket:   b.Name,
		Key:      key,
		UploadId: domain.RandomId(),
	})
}

func handleUploadPart(c *srv.Context) *srv.Response {
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		return responseFromError(ec.InvalidArgument)
	}
	partNumber, r := c.IntQuery("part-number")
	if r != nil {
		return r
	}

	slog.Info("upload part", "upload", uploadId, "part", partNumber)

	return srv.Respond().ETag(domain.NewEtag())
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

func handleCompleteMultipartUpload(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	key, r := contextGetObjectKey(c)
	if r != nil {
		return r
	}
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		return responseFromError(ec.InvalidArgument)
	}
	var req CompleteMultipartUploadRequest
	if r := c.BindJSON(&req); r != nil {
		return r
	}

	slog.Info("complete multipart upload", "bucket", b.Name, "key", key, "upload", uploadId)

	return srv.Respond().Json(CompleteMultipartUploadResult{
		Bucket: b.Name,
		Key:    key,
		ETag:   domain.RandomId(),
	})
}

func handleAbortMultipartUpload(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	key, r := contextGetObjectKey(c)
	if r != nil {
		return r
	}
	uploadId := c.Query(queryUploadId)
	if uploadId == "" {
		return responseFromError(ec.InvalidArgument)
	}

	slog.Info("abort multipart upload", "bucket", b.Name, "key", key, "upload", uploadId)

	return srv.Respond().NoContent()
}
