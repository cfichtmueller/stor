// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log"

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

type AddArchiveEntriesRequest struct {
	Entries []archive.Entry `json:"entries"`
}

func handleAddArchiveEntries(c jug.Context) {
	archiveId, ok := archiveFilter(c)
	if !ok {
		return
	}
	var req AddArchiveEntriesRequest
	if !c.MustBindJSON(&req) {
		return
	}

	if err := archive.AddEntries(c, archive.AddEntriesCommand{
		ArchiveId: archiveId,
		Entries:   req.Entries,
	}); err != nil {
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
	archiveId, ok := archiveFilter(c)
	if !ok {
		return
	}

	if err := archive.Complete(c, archiveId); err != nil {
		handleError(c, err)
		return
	}

	c.RespondNoContent()
}

func handleAbortArchive(c jug.Context) {
	archiveId, ok := archiveFilter(c)
	if !ok {
		return
	}

	log.Printf("abort archive %s", archiveId)

	c.RespondNoContent()
}
