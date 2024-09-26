// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package archive

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cfichtmueller/stor/internal/bus"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/ec"
)

var (
	//EventCompleted indicates that an archive has been completed. Data is a CompletedEvent
	EventCompleted       = "archive.completed"
	StatePending         = "pending"
	StateProcessing      = "processing"
	StateComplete        = "complete"
	StateFailed          = "failed"
	TypeZip              = "zip"
	archiveFields        = "id, bucket, key, type, state"
	createStmt           *sql.Stmt
	findOneStmt          *sql.Stmt
	findOneWithStateStmt *sql.Stmt
	existsStmt           *sql.Stmt
	updateStmt           *sql.Stmt
	deleteStmt           *sql.Stmt
	insertEntryStmt      *sql.Stmt
	findEntriesStmt      *sql.Stmt
	deleteEntriesStmt    *sql.Stmt
	completeMutex        sync.Mutex
	finishFlag           = true
)

type CompletedEvent struct {
	Bucket    string
	Key       string
	ArchiveId string
}

func Configure() {
	db.RunMigration("create_archive_table", `CREATE TABLE archives (
	id CHAR(32) PRIMARY KEY,
	bucket CHAR(64) NOT NULL,
	key TEXT NOT NULL,
	type CHAR(6) NOT NULL,
	state CHAR(64) NOT NULL, 
	is_deleted INTEGER NOT NULL
	)`)

	db.RunMigration("create_archive_entries_table", `CREATE TABLE archive_entries (
	id CHAR(32) PRIMARY KEY,
	archive CHAR(32) NOT NULL,
	key TEXT NOT NULL,
	name TEXT NOT NULL
	)`)

	db.RunMigration("create_archive_entries_index", `CREATE INDEX idx_archive_entries ON archive_entries (archive, name)`)

	createStmt = db.Prepare("INSERT INTO archives (id, bucket, key, type, state, is_deleted) VALUES ($1, $2, $3, $4, $5, false)")
	findOneStmt = db.Prepare("SELECT " + archiveFields + " FROM archives WHERE id = $1 AND is_deleted = $2")
	findOneWithStateStmt = db.Prepare("SELECT " + archiveFields + " FROM archives WHERE state = $1 AND is_deleted = false LIMIT 1")
	existsStmt = db.Prepare("SELECT COUNT(*) FROM archives WHERE id = $1 AND bucket = $2 AND key = $3 AND is_deleted = $4")
	updateStmt = db.Prepare("UPDATE archives SET state = $1, is_deleted = $2 WHERE id = $3")
	deleteStmt = db.Prepare("DELETE FROM archives WHERE id = $1")
	insertEntryStmt = db.Prepare("INSERT INTO archive_entries (id, archive, key, name) VALUES ($1, $2, $3, $4)")
	findEntriesStmt = db.Prepare("SELECT key, name FROM archive_entries WHERE archive = $1 AND name > $2 LIMIT $3")
	deleteEntriesStmt = db.Prepare("DELETE FROM archive_entries WHERE archive = $1")

	go worker()
}

type Archive struct {
	ID     string
	Bucket string
	Key    string
	Type   string
	State  string
}

type CreateCommand struct {
	Bucket string
	Key    string
	Type   string
}

func Create(ctx context.Context, cmd CreateCommand) (string, error) {
	if cmd.Type != TypeZip {
		return "", ec.InvalidArgument
	}
	id := domain.RandomId()
	if _, err := createStmt.ExecContext(ctx, id, cmd.Bucket, cmd.Key, cmd.Type, StatePending); err != nil {
		return "", fmt.Errorf("unable to create archive record: %v", err)
	}
	return id, nil
}

func Exists(ctx context.Context, bucket, key, id string) (bool, error) {
	var count int
	if err := existsStmt.QueryRowContext(ctx, id, bucket, key, false).Scan(&count); err != nil {
		return false, fmt.Errorf("unable to query archives: %v", err)
	}
	return count > 0, nil
}

func FindOne(ctx context.Context, bucket, key, id string) (*Archive, error) {
	arch, err := scanRow(findOneStmt.QueryRowContext(ctx, id, false))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.NoSuchArchive
		}
		return nil, fmt.Errorf("unable to find archive: %v", err)
	}
	if arch.Bucket != bucket || arch.Key != key {
		return nil, ec.NoSuchArchive
	}
	return arch, nil
}

type Entry struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func AddEntries(ctx context.Context, a *Archive, entries []Entry) error {
	//TODO: should be bulk insert
	for _, e := range entries {
		if _, err := insertEntryStmt.ExecContext(ctx, domain.RandomId(), a.ID, e.Key, e.Name); err != nil {
			return fmt.Errorf("unable to insert entry record: %v", err)
		}
	}
	return nil
}

type CompleteResult struct {
	Bucket string
	Key    string
	ETag   string
}

func Complete(ctx context.Context, a *Archive) error {
	completeMutex.Lock()
	defer completeMutex.Unlock()

	if a.State != StatePending {
		return ec.ArchiveNotPending
	}

	if _, err := updateStmt.ExecContext(ctx, StateProcessing, false, a.ID); err != nil {
		return fmt.Errorf("unable to update archive record: %v", err)
	}
	finishFlag = true

	return nil
}

func Abort(ctx context.Context, a *Archive) error {
	if a.State != StatePending && a.State != StateFailed {
		return ec.ArchiveNotAbortable
	}
	return delete(ctx, a.ID)
}

func delete(ctx context.Context, id string) error {
	if _, err := deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("unable to delete archive record: %v", err)
	}

	if _, err := deleteEntriesStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("unable to delete archive entries: %v", err)
	}
	return nil
}

func worker() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		finishArchives()
	}
}

func finishArchives() {
	if !finishFlag {
		return
	}
	ctx := context.Background()
	a, err := scanRow(findOneWithStateStmt.QueryRowContext(ctx, StateProcessing))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			finishFlag = false
			return
		}
		log.Printf("unable to query pending archives: %v", err)
		return
	}

	if err := finishArchive(ctx, a); err != nil {
		log.Printf("unable to finish archive: %v", err)
		return
	}

	finishFlag = true
}

func finishArchive(ctx context.Context, arch *Archive) error {
	s := NewStats()
	chunkWriter, err := chunk.NewWriter()
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(chunkWriter)

	startAfter := ""
	for {
		rows, err := findEntriesStmt.QueryContext(ctx, arch.ID, startAfter, 1000)
		if err != nil {
			return fmt.Errorf("unable to list archive entries: %v", err)
		}
		entries := make([]Entry, 0)
		for rows.Next() {
			var entry Entry
			if err := rows.Scan(&entry.Key, &entry.Name); err != nil {
				return fmt.Errorf("unable to scan entry row: %v", err)
			}
			entries = append(entries, entry)
		}
		if len(entries) == 0 {
			break
		}

		for _, e := range entries {
			o, err := object.FindOne(ctx, arch.Bucket, e.Key, false)
			if err != nil {
				if errors.Is(err, ec.NoSuchKey) {
					failArchive(ctx, arch.ID)
					return nil
				}
				return err
			}

			writer, err := zipWriter.Create(e.Name)
			if err != nil {
				return fmt.Errorf("unable to create zip entry: %v", err)
			}
			if err := object.Write(ctx, o, writer); err != nil {
				return fmt.Errorf("unable to write object: %v", err)
			}
			startAfter = e.Name
			s.AddBytes(o.Size)
		}
		s.AddFiles(int64(len(entries)))
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("unable to close zip writer: %v", err)
	}

	chunkId, err := chunkWriter.Commit(ctx)
	if err != nil {
		return err
	}

	existing, err := object.FindOne(ctx, arch.Bucket, arch.Key, false)
	if err != nil && !errors.Is(err, ec.NoSuchKey) {
		return err
	}
	if err == nil {
		if err := object.Delete(ctx, existing); err != nil {
			return err
		}
	}

	if _, err := object.CreateWithChunk(ctx, arch.Bucket, chunkId, object.CreateCommand{
		Key:         arch.Key,
		ContentType: "application/zip",
		Size:        chunkWriter.Size(),
	}); err != nil {
		return err
	}

	if err := delete(ctx, arch.ID); err != nil {
		fmt.Errorf("unable to delete archive: %v", err)
	}

	log.Printf("Finished archive: %s", s.Summary())

	bus.Publish(EventCompleted, CompletedEvent{
		Bucket:    arch.Bucket,
		Key:       arch.Key,
		ArchiveId: arch.ID,
	})

	return nil
}

func failArchive(ctx context.Context, id string) {
	if _, err := updateStmt.ExecContext(ctx, StateFailed, false, id); err != nil {
		log.Printf("unable to fail archive %s: %v", id, err)
	}
}

func scanRow(row *sql.Row) (*Archive, error) {
	a := Archive{}
	if err := row.Scan(&a.ID, &a.Bucket, &a.Key, &a.Type, &a.State); err != nil {
		return nil, err
	}

	return &a, nil
}
