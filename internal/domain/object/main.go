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
	ContentType string
	Size        uint64
	CreatedAt   time.Time
}

var (
	createStmt             *sql.Stmt
	listStmt               *sql.Stmt
	findOneStmt            *sql.Stmt
	existsStmt             *sql.Stmt
	deleteStmt             *sql.Stmt
	addObjectChunkStmt     *sql.Stmt
	findObjectChunksStmt   *sql.Stmt
	deleteObjectChunksStmt *sql.Stmt
	countStmt              *sql.Stmt
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
		"SELECT id, bucket, key, content_type, size, created_at FROM objects WHERE bucket = $1 AND key > $2 ORDER BY key LIMIT $3",
		"SELECT id, bucket, key, content_type, size, created_at FROM objects WHERE bucket = $1 AND key = $2 LIMIT 1",
		"SELECT COUNT(*) as count FROM objects WHERE bucket = $1 AND key = $2",
		"DELETE FROM objects WHERE id = $1",
		"INSERT INTO object_chunks (object, chunk, seq) VALUES ($1, $2, $3)",
		"SELECT chunk FROM object_chunks WHERE object = $1 ORDER BY seq",
		"DELETE FROM object_chunks WHERE object = $1",
		"SELECT COUNT(*) FROM objects WHERE bucket = $1 AND key > $2",
	)

	createStmt = s[0]
	listStmt = s[1]
	findOneStmt = s[2]
	existsStmt = s[3]
	deleteStmt = s[4]
	addObjectChunkStmt = s[5]
	findObjectChunksStmt = s[6]
	deleteObjectChunksStmt = s[7]
	countStmt = s[8]
}

func List(ctx context.Context, bucketName, startAfter string, limit int) ([]*Object, error) {
	return decodeRows(listStmt.QueryContext(ctx, bucketName, startAfter, limit))
}

func Count(ctx context.Context, bucketName, startAfter string) (int, error) {
	var count int
	if err := countStmt.QueryRowContext(ctx, bucketName, startAfter).Scan(&count); err != nil {
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
			return nil, ec.NoSuchKey
		}
		return nil, fmt.Errorf("unable to find object record: %v", err)
	}
	return &o, nil
}

func FindMany(ctx context.Context, bucketName string, keys []string) ([]*Object, error) {
	query := strings.Builder{}
	query.WriteString("SELECT id, bucket, key, content_type, size, created_at FROM objects WHERE bucket = $1 AND key IN (")
	first := true
	for i := range keys {
		if first {
			first = false
		} else {
			query.WriteString(", ")
		}
		query.WriteString("$" + strconv.Itoa(i+2))
	}
	query.WriteString(")")
	stmt, err := db.PrepareOne(query.String())
	if err != nil {
		return nil, err
	}
	args := make([]any, 0, len(keys)+1)
	args = append(args, bucketName)
	for _, k := range keys {
		args = append(args, k)
	}
	return decodeRows(stmt.QueryContext(ctx, args...))
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

	if _, err := addObjectChunkStmt.ExecContext(ctx, objectId, chunkId, 1); err != nil {
		return fmt.Errorf("unable to persist object chunk record: %v", err)
	}

	return nil
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
	chunkIds, err := findObjectChunks(ctx, o.ID)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Delete(ctx, chunkId); err != nil {
			return err
		}
	}

	if _, err := deleteObjectChunksStmt.ExecContext(ctx, o.ID); err != nil {
		return fmt.Errorf("unable to delete chunk links: %v", err)
	}

	if _, err := deleteStmt.ExecContext(ctx, o.ID); err != nil {
		return fmt.Errorf("unable to delete object: %v", err)
	}

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
