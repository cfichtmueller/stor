// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
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

	log.Printf("Create archive in `%s/%s` -> %s", b.Name, key, t)

	c.RespondOk(CreateArchiveResult{
		Bucket:    b.Name,
		Key:       key,
		ArchiveId: domain.RandomId(),
	})
}

type ArchiveEntry struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type AddArchiveEntriesRequest struct {
	Entries []ArchiveEntry `json:"entries"`
}

func handleAddArchiveEntries(c jug.Context) {
	archiveId := c.Query(queryArchiveId)
	if archiveId == "" {
		handleError(c, ec.InvalidArgument)
		return
	}

	var req AddArchiveEntriesRequest
	if !c.MustBindJSON(&req) {
		return
	}

	for _, e := range req.Entries {
		log.Printf("add entry %s -> %s to archive %s", e.Key, e.Name, archiveId)
	}

	c.Status(200)
}

type CompleteArchiveResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	ETag   string `json:"etag"`
}

func handleCompleteArchive(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	archiveId := c.Query(queryArchiveId)
	if archiveId == "" {
		handleError(c, ec.InvalidArgument)
		return
	}

	log.Printf("complete archive in `%s/%s` -> %s", b.Name, key, archiveId)

	c.RespondOk(CompleteArchiveResult{
		Bucket: b.Name,
		Key:    key,
		ETag:   domain.RandomId(),
	})
}

func handleAbortArchive(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)
	archiveId := c.Query(queryArchiveId)
	if archiveId == "" {
		handleError(c, ec.InvalidArgument)
		return
	}

	log.Printf("abort archive in `%s/%s` -> %s", b.Name, key, archiveId)

	c.RespondNoContent()
}
