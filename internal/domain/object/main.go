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
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
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
	ContentType string
	Size        uint64
	CreatedAt   time.Time
}

var (
	ErrNotFound      = jug.NewNotFoundError("object not found")
	createStmt       *sql.Stmt
	findManyStmt     *sql.Stmt
	findOneStmt      *sql.Stmt
	existsStmt       *sql.Stmt
	deleteStmt       *sql.Stmt
	addChunkStmt     *sql.Stmt
	findChunksStmt   *sql.Stmt
	deleteChunksStmt *sql.Stmt
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

	s := db.Prepare(
		"INSERT INTO objects (id, bucket, key, content_type, size, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		"SELECT id, bucket, key, content_type, size, created_at FROM objects WHERE bucket = $1 ORDER BY key",
		"SELECT id, bucket, key, content_type, size, created_at FROM objects WHERE bucket = $1 AND key = $2 LIMIT 1",
		"SELECT COUNT(*) as count FROM objects WHERE bucket = $1 AND key = $2",
		"DELETE FROM objects WHERE id = $1",
		"INSERT INTO object_chunks (object, chunk, seq) VALUES ($1, $2, $3)",
		"SELECT chunk FROM object_chunks WHERE object = $1 ORDER BY seq",
		"DELETE FROM object_chunks WHERE object = $1",
	)

	createStmt = s[0]
	findManyStmt = s[1]
	findOneStmt = s[2]
	existsStmt = s[3]
	deleteStmt = s[4]
	addChunkStmt = s[5]
	findChunksStmt = s[6]
	deleteChunksStmt = s[7]
}

func List(ctx context.Context, bucketName string) ([]*Object, error) {
	return decodeRows(findManyStmt.QueryContext(ctx, bucketName))
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
			&o.ContentType,
			&o.Size,
			&o.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to decode object record: %v", err)
		}
		objects = append(objects, &o)
	}
	return objects, nil
}

func FindOne(ctx context.Context, bucketName, key string) (*Object, error) {
	var o Object
	if err := findOneStmt.QueryRowContext(ctx, bucketName, key).Scan(
		&o.ID,
		&o.Bucket,
		&o.Key,
		&o.ContentType,
		&o.Size,
		&o.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unable to find object record: %v", err)
	}
	return &o, nil
}

func Exists(ctx context.Context, bucketName, key string) (bool, error) {
	var count int
	if err := existsStmt.QueryRowContext(ctx, bucketName, key).Scan(&count); err != nil {
		return false, fmt.Errorf("unable to count objects: %v", err)
	}
	return count > 0, nil
}

func Create(ctx context.Context, bucketId string, cmd CreateCommand) error {
	objectId := domain.RandomId()
	chunkId, err := chunk.Create(ctx, cmd.Data)
	if err != nil {
		return err
	}

	if _, err := createStmt.ExecContext(ctx, objectId, bucketId, cmd.Key, cmd.ContentType, len(cmd.Data), time.Now()); err != nil {
		return fmt.Errorf("unable to persist object record: %v", err)
	}

	if _, err := addChunkStmt.ExecContext(ctx, objectId, chunkId, 1); err != nil {
		return fmt.Errorf("unable to persist object chunk record: %v", err)
	}

	return nil
}

func Write(ctx context.Context, o *Object, w io.Writer) error {
	chunkIds, err := findChunks(ctx, o.ID)
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
	chunkIds, err := findChunks(ctx, o.ID)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Delete(ctx, chunkId); err != nil {
			return err
		}
	}

	if _, err := deleteChunksStmt.ExecContext(ctx, o.ID); err != nil {
		return fmt.Errorf("unable to delete chunk links: %v", err)
	}

	if _, err := deleteStmt.ExecContext(ctx, o.ID); err != nil {
		return fmt.Errorf("unable to delete object: %v", err)
	}

	return nil
}

func findChunks(ctx context.Context, objectId string) ([]string, error) {
	//TODO: when the number of chunks becomes large, this needs to "cursor" its way through
	rows, err := findChunksStmt.QueryContext(ctx, objectId)
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
