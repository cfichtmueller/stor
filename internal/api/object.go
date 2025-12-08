// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"io"
	"time"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
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

func handleObjectHead(c *srv.Context) *srv.Response {
	_, ok, err := authenticateApiKey(c)
	if err != nil {
		return responseFromError(err)
	}
	if !ok {
		if r := mustAuthenticateNonce(c); r != nil {
			return r
		}
	}
	b, r := mustGetBucket(c)
	if r != nil {
		return r
	}
	o, r := mustGetObject(c, b)
	if r != nil {
		return r
	}
	return srv.Respond().
		ContentLength(int64(o.Size)).
		ContentType(o.ContentType)
}

func handleObjectGet(c *srv.Context) *srv.Response {
	query := c.Request().URL.Query()
	if query.Has(queryArchiveId) {
		return handleGetArchive(c)
	}
	return handleGetObject(c)
}

func handleGetObject(c *srv.Context) *srv.Response {
	_, ok, err := authenticateApiKey(c)
	if err != nil {
		return responseFromError(err)
	}
	if !ok {
		if r := mustAuthenticateNonce(c); r != nil {
			return r
		}
	}
	b, r := mustGetBucket(c)
	if r != nil {
		return r
	}
	o, r := mustGetObject(c, b)
	if r != nil {
		return r
	}
	return srv.Respond().
		ContentLength(int64(o.Size)).
		BodyFn(o.ContentType, func(w io.Writer) error {
			return object.Write(c, o, w)
		})
}

func handleObjectPost(c *srv.Context) *srv.Response {
	if r := mustAuthenticateApiKey(c); r != nil {
		return r
	}
	b, r := mustGetBucket(c)
	if r != nil {
		return r
	}
	contextSetBucket(c, b)

	if c.HasQuery(queryArchives) {
		return handleCreateArchive(c)
	} else if c.Query(queryArchiveId) != "" {
		return handleCompleteArchive(c)
	} else if c.HasQuery(queryNonces) {
		return handleCreateNonce(c)
	} else if c.HasQuery(queryUploads) {
		return handleCreateMultipartUpload(c)
	} else if c.Query(queryUploadId) != "" {
		return handleCompleteMultipartUpload(c)
	}
	return srv.Respond().MethodNotAllowed()
}

func handleObjectPut(c *srv.Context) *srv.Response {
	if c.Query(queryArchiveId) != "" {
		return handleAddArchiveEntries(c)
	} else if c.Query(queryUploadId) != "" {
		return handleUploadPart(c)
	}
	return handleCreateOrUpdateObject(c)
}

func handleObjectDelete(c *srv.Context) *srv.Response {
	if c.Query(queryArchiveId) != "" {
		return handleAbortArchive(c)
	} else if c.Query(queryUploadId) != "" {
		return handleAbortMultipartUpload(c)
	}
	return handleDeleteObject(c)
}

func handleCreateOrUpdateObject(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	key, r := contextGetObjectKey(c)
	if r != nil {
		return r
	}
	copySource := c.Header("Stor-Copy-Source")

	if copySource != "" {
		return createOrUpdateObjectFromCopySource(c, b, key, copySource)
	}

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	d, err := c.GetRawData()
	if err != nil {
		if errors.Is(err, srv.ErrNoBody) {
			return srv.Respond().BadRequest(srv.ErrorDto{
				Code:    "request_body_missing",
				Message: "Request body is missing",
			})
		}
		return responseFromError(err)
	}

	exists, err := object.Exists(c, b.Name, key)
	if err != nil {
		return responseFromError(err)
	}

	if exists {
		existing, err := object.FindOne(c, b.Name, key, false)
		if err != nil {
			return responseFromError(err)
		}
		updated, err := uc.UpdateObjectWithData(c, b, existing, object.UpdateCommand{
			ContentType: contentType,
			Data:        d,
		})
		if err != nil {
			return responseFromError(err)
		}
		return srv.Respond().NoContent().ETag(updated.ETag)
	}

	created, err := uc.CreateObjectFromData(c, b, object.CreateCommand{
		Key:         key,
		ContentType: contentType,
		Data:        d,
	})
	if err != nil {
		return responseFromError(err)
	}
	return srv.Respond().NoContent().ETag(created.ETag)
}

func createOrUpdateObjectFromCopySource(c *srv.Context, b *bucket.Bucket, key, copySource string) *srv.Response {
	src, err := object.FindOne(c, b.Name, copySource, false)
	if err != nil {
		return responseFromError(err)
	}
	exists, err := object.Exists(c, b.Name, key)
	if err != nil {
		return responseFromError(err)
	}
	if exists {
		existing, err := object.FindOne(c, b.Name, key, false)
		if err != nil {
			return responseFromError(err)
		}
		updated, err := uc.UpdateObjectFromCopy(c, b, src, existing)
		if err != nil {
			return responseFromError(err)
		}
		return srv.Respond().NoContent().ETag(updated.ETag)
	}
	created, err := uc.CreateObjectFromCopy(c, b, src, key)
	if err != nil {
		return responseFromError(err)
	}
	return srv.Respond().NoContent().ETag(created.ETag)
}

func handleDeleteObject(c *srv.Context) *srv.Response {
	o, r := objectFilter(c)
	if r != nil {
		return r
	}
	b := contextGetBucket(c)
	if err := uc.DeleteObject(c, b, o); err != nil {
		return responseFromError(err)
	}

	return srv.Respond().NoContent()
}
