// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package object

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/ec"
)

type CreateCommand struct {
	Key         string
	ContentType string
	Data        []byte
}

type Object struct {
	ID          string
	Bucket      string
	Key         string
	ETag        string
	ContentType string
	Size        uint64
	CreatedAt   time.Time
	Deleted     bool
}

var (
	maxPurgeTime = 400 * time.Millisecond

	createStmt             *sql.Stmt
	listStmt               *sql.Stmt
	findOneStmt            *sql.Stmt
	existsStmt             *sql.Stmt
	updateStmt             *sql.Stmt
	deleteStmt             *sql.Stmt
	findDeletedStmt        *sql.Stmt
	addObjectChunkStmt     *sql.Stmt
	findObjectChunksStmt   *sql.Stmt
	deleteObjectChunksStmt *sql.Stmt
	countStmt              *sql.Stmt
	purgeFlag              = true
)

func Configure() {
	if err := config.Mkdir("chunks"); err != nil {
		log.Fatalf("unable to create chunk directory: %v", err)
	}

	db.RunMigration("create_object_table", `CREATE TABLE objects(
		id CHAR(32) PRIMARY KEY,
		bucket CHAR(64) NOT NULL,
		key TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INT NOT NULL,
		created_at DATETIME NOT NULL
	)
	`)

	db.RunMigration("create_object_chunks_table", `CREATE TABLE object_chunks(
		object CHAR(32),
		chunk CHAR(64),
		seq INT,
		PRIMARY KEY (object, chunk, seq)
	)`)

	db.RunMigration("add_object_deleted_flag", `ALTER TABLE objects ADD COLUMN is_deleted INTEGER`)
	db.RunMigration("add_object_key_index", `CREATE INDEX idx_objects_key_bucket_deleted ON objects (key, bucket, is_deleted)`)
	db.RunMigration("add_objectchunk_key_index", `CREATE INDEX idx_objectchunk_key ON object_chunks (object)`)
	db.RunMigration("add_object_etag_1", `ALTER TABLE objects ADD COLUMN etag CHAR(64)`)

	createStmt = db.Prepare("INSERT INTO objects (id, bucket, key, etag, content_type, size, created_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, false)")
	listStmt = db.Prepare("SELECT id, bucket, key, etag, content_type, size, created_at, is_deleted FROM objects WHERE bucket = $1 AND key > $2 AND is_deleted = $3 ORDER BY key LIMIT $4")
	findOneStmt = db.Prepare("SELECT id, bucket, key, etag, content_type, size, created_at, is_deleted FROM objects WHERE bucket = $1 AND key = $2 AND is_deleted = $3 LIMIT 1")
	existsStmt = db.Prepare("SELECT COUNT(*) as count FROM objects WHERE bucket = $1 AND key = $2 AND is_deleted = $3")
	updateStmt = db.Prepare("UPDATE objects SET is_deleted = $1 WHERE id = $2")
	deleteStmt = db.Prepare("DELETE FROM objects WHERE id = $1")
	findDeletedStmt = db.Prepare("SELECT id FROM objects WHERE is_deleted = true LIMIT 1000")
	addObjectChunkStmt = db.Prepare("INSERT INTO object_chunks (object, chunk, seq) VALUES ($1, $2, $3)")
	findObjectChunksStmt = db.Prepare("SELECT chunk FROM object_chunks WHERE object = $1 ORDER BY seq")
	deleteObjectChunksStmt = db.Prepare("DELETE FROM object_chunks WHERE object = $1")
	countStmt = db.Prepare("SELECT COUNT(*) FROM objects WHERE bucket = $1 AND key > $2 AND is_deleted  = $3")

	go worker()
}

func List(ctx context.Context, bucketName, startAfter string, limit int) ([]*Object, error) {
	return decodeRows(listStmt.QueryContext(ctx, bucketName, startAfter, false, limit))
}

func Count(ctx context.Context, bucketName, startAfter string) (int, error) {
	var count int
	if err := countStmt.QueryRowContext(ctx, bucketName, startAfter, false).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to count objects: %v", err)
	}
	return count, nil
}

func decodeRows(rows *sql.Rows, err error) ([]*Object, error) {
	if err != nil {
		return nil, fmt.Errorf("unable to find object records: %v", err)
	}
	objects := make([]*Object, 0)
	for rows.Next() {
		var o Object
		if err := rows.Scan(
			&o.ID,
			&o.Bucket,
			&o.Key,
			&o.ETag,
			&o.ContentType,
			&o.Size,
			&o.CreatedAt,
			&o.Deleted,
		); err != nil {
			return nil, fmt.Errorf("unable to decode object record: %v", err)
		}
		objects = append(objects, &o)
	}
	return objects, nil
}

func FindOne(ctx context.Context, bucketName, key string, deleted bool) (*Object, error) {
	var o Object
	if err := findOneStmt.QueryRowContext(ctx, bucketName, key, deleted).Scan(
		&o.ID,
		&o.Bucket,
		&o.Key,
		&o.ETag,
		&o.ContentType,
		&o.Size,
		&o.CreatedAt,
		&o.Deleted,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.NoSuchKey
		}
		return nil, fmt.Errorf("unable to find object record: %v", err)
	}
	return &o, nil
}

func FindMany(ctx context.Context, bucketName string, keys []string, deleted bool) ([]*Object, error) {
	query := strings.Builder{}
	query.WriteString("SELECT id, bucket, key, content_type, size, created_at, is_deleted FROM objects WHERE bucket = $1 AND is_deleted = $2 AND key IN (")
	first := true
	for i := range keys {
		if first {
			first = false
		} else {
			query.WriteString(", ")
		}
		query.WriteString("$" + strconv.Itoa(i+3))
	}
	query.WriteString(")")
	stmt, err := db.PrepareOne(query.String())
	if err != nil {
		return nil, err
	}
	args := make([]any, 0, len(keys)+3)
	args = append(args, bucketName)
	args = append(args, deleted)
	for _, k := range keys {
		args = append(args, k)
	}
	return decodeRows(stmt.QueryContext(ctx, args...))
}

func Exists(ctx context.Context, bucketName, key string) (bool, error) {
	var count int
	if err := existsStmt.QueryRowContext(ctx, bucketName, key, false).Scan(&count); err != nil {
		return false, fmt.Errorf("unable to count objects: %v", err)
	}
	return count > 0, nil
}

func Create(ctx context.Context, bucketId string, cmd CreateCommand) (*Object, error) {
	chunkId, err := chunk.Create(ctx, cmd.Data)
	if err != nil {
		return nil, err
	}

	o := Object{
		ID:          domain.RandomId(),
		Bucket:      bucketId,
		Key:         cmd.Key,
		ETag:        domain.NewEtag(),
		ContentType: cmd.ContentType,
		Size:        uint64(len(cmd.Data)),
		CreatedAt:   time.Now(),
	}

	if _, err := createStmt.ExecContext(ctx, &o.ID, &o.Bucket, &o.Key, &o.ETag, &o.ContentType, &o.Size, &o.CreatedAt); err != nil {
		return nil, fmt.Errorf("unable to persist object record: %v", err)
	}

	if _, err := addObjectChunkStmt.ExecContext(ctx, o.ID, chunkId, 1); err != nil {
		return nil, fmt.Errorf("unable to persist object chunk record: %v", err)
	}

	return &o, nil
}

func Write(ctx context.Context, o *Object, w io.Writer) error {
	chunkIds, err := findObjectChunks(ctx, o.ID)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Write(chunkId, w); err != nil {
			return err
		}
	}
	return nil
}

func Delete(ctx context.Context, o *Object) error {
	o.Deleted = true
	if _, err := updateStmt.ExecContext(ctx, true, o.ID); err != nil {
		return fmt.Errorf("unable to update object record: %v", err)
	}

	purgeFlag = true

	return nil
}

func findObjectChunks(ctx context.Context, objectId string) ([]string, error) {
	//TODO: when the number of chunks becomes large, this needs to "cursor" its way through
	rows, err := findObjectChunksStmt.QueryContext(ctx, objectId)
	if err != nil {
		return nil, fmt.Errorf("unable to find chunks: %v", err)
	}
	ids := make([]string, 0)
	for rows.Next() {
		var chunkId string
		if err := rows.Scan(&chunkId); err != nil {
			return nil, fmt.Errorf("unable to decode chunk: %v", err)
		}
		ids = append(ids, chunkId)
	}
	return ids, nil
}

func worker() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		purgeObjects()
	}
}

func purgeObjects() {
	if !purgeFlag {
		return
	}
	start := time.Now()
	ctx := context.Background()
	rows, err := findDeletedStmt.QueryContext(ctx)
	if err != nil {
		log.Printf("unable to purge objects: %v", err)
		return
	}
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("unable to purge objects: %v", err)
			return
		}
		ids = append(ids, id)
	}
	for i, id := range ids {
		if err := purgeObject(ctx, id); err != nil {
			log.Printf("unable to purge objects: %v", err)
			return
		}
		if time.Since(start) > maxPurgeTime {
			if i > 1 {
				log.Printf("purged %d objects", i-1)
			}
			return
		}
	}
	if len(ids) > 0 {
		log.Printf("purged %d objects", len(ids))
	}
	purgeFlag = false
}

func purgeObject(ctx context.Context, objectId string) error {
	chunkIds, err := findObjectChunks(ctx, objectId)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Delete(ctx, chunkId); err != nil {
			return err
		}
	}

	if _, err := deleteObjectChunksStmt.ExecContext(ctx, objectId); err != nil {
		return fmt.Errorf("unable to delete chunk links: %v", err)
	}

	if _, err := deleteStmt.ExecContext(ctx, objectId); err != nil {
		return fmt.Errorf("unable to delete object: %v", err)
	}

	return nil
}
