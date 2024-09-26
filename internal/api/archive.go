// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/archive"
)

type CreateArchiveResult struct {
	Bucket    string `json:"bucket"`
	Key       string `json:"key"`
	ArchiveId string `json:"archiveId"`
}

func handleCreateArchive(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	t := c.Query("type")

	id, err := archive.Create(c, archive.CreateCommand{
		Bucket: b.Name,
		Key:    key,
		Type:   t,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.RespondOk(CreateArchiveResult{
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

func handleGetArchive(c jug.Context) {
	if !mustAuthenticateApiKey(c) {
		return
	}
	b, ok := mustGetBucket(c)
	if !ok {
		return
	}
	contextSetBucket(c, b)
	arch, ok := archiveFilter(c)
	if !ok {
		return
	}
	c.RespondOk(ArchiveResponse{
		ID:    arch.ID,
		State: arch.State,
		Type:  arch.Type,
	})
}

type AddArchiveEntriesRequest struct {
	Entries []archive.Entry `json:"entries"`
}

func handleAddArchiveEntries(c jug.Context) {
	arch, ok := archiveFilter(c)
	if !ok {
		return
	}
	var req AddArchiveEntriesRequest
	if !c.MustBindJSON(&req) {
		return
	}

	if err := archive.AddEntries(c, arch, req.Entries); err != nil {
		handleError(c, err)
		return
	}

	c.Status(200)
}

type CompleteArchiveResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	ETag   string `json:"etag"`
}

func handleCompleteArchive(c jug.Context) {
	arch, ok := archiveFilter(c)
	if !ok {
		return
	}

	if err := archive.Complete(c, arch); err != nil {
		handleError(c, err)
		return
	}

	c.RespondNoContent()
}

func handleAbortArchive(c jug.Context) {
	arch, ok := archiveFilter(c)
	if !ok {
		return
	}

	if err := archive.Abort(c, arch); err != nil {
		handleError(c, err)
	}

	c.RespondNoContent()
}
