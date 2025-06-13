// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/archive"
)

type CreateArchiveResult struct {
	Bucket    string `json:"bucket"`
	Key       string `json:"key"`
	ArchiveId string `json:"archiveId"`
}

func handleCreateArchive(c *srv.Context) *srv.Response {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	t := c.Query("type")

	id, err := archive.Create(c, archive.CreateCommand{
		Bucket: b.Name,
		Key:    key,
		Type:   t,
	})
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().Json(CreateArchiveResult{
		Bucket:    b.Name,
		Key:       key,
		ArchiveId: id,
	})
}

type ArchiveResponse struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Type  string `json:"type"`
}

func handleGetArchive(c *srv.Context) *srv.Response {
	if r := mustAuthenticateApiKey(c); r != nil {
		return r
	}
	b, r := mustGetBucket(c)
	if r != nil {
		return r
	}
	contextSetBucket(c, b)
	arch, r := archiveFilter(c)
	if r != nil {
		return r
	}
	return srv.Respond().Json(ArchiveResponse{
		ID:    arch.ID,
		State: arch.State,
		Type:  arch.Type,
	})
}

type AddArchiveEntriesRequest struct {
	Entries []archive.Entry `json:"entries"`
}

func handleAddArchiveEntries(c *srv.Context) *srv.Response {
	arch, r := archiveFilter(c)
	if r != nil {
		return r
	}
	var req AddArchiveEntriesRequest
	if r := c.BindJSON(&req); r != nil {
		return r
	}

	if err := archive.AddEntries(c, arch, req.Entries); err != nil {
		return responseFromError(err)
	}

	return srv.Respond()
}

type CompleteArchiveResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	ETag   string `json:"etag"`
}

func handleCompleteArchive(c *srv.Context) *srv.Response {
	arch, r := archiveFilter(c)
	if r != nil {
		return r
	}

	if err := archive.Complete(c, arch); err != nil {
		return responseFromError(err)
	}

	return srv.Respond().NoContent()
}

func handleAbortArchive(c *srv.Context) *srv.Response {
	arch, r := archiveFilter(c)
	if r != nil {
		return r
	}

	if err := archive.Abort(c, arch); err != nil {
		return responseFromError(err)
	}

	return srv.Respond().NoContent()
}
