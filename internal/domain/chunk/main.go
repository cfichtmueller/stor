// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chunk

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"sync"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
)

type Chunk struct {
	ID         string
	Size       uint64
	References uint64
}

type Stats struct {
	Count     uint64
	TotalSize uint64
}

var (
	ErrNotFound = fmt.Errorf("chunk not found")
	createStmt  *sql.Stmt
	findOneStmt *sql.Stmt
	updateStmt  *sql.Stmt
	deleteStmt  *sql.Stmt
	statsStmt   *sql.Stmt
	writeMutex  sync.Mutex
)

func Configure() {
	if err := os.Mkdir(path.Join(config.DataDir, "chunks"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("unable to create chunk directory: %v", err)
		}
	}

	db.RunMigration("create_chunks_table", `CREATE TABLE chunks(
		id CHAR(64) PRIMARY KEY,
		size INT NOT NULL,
		rc INT NOT NULL
	)`)

	createStmt = db.Prepare("INSERT INTO chunks (id, size, rc) VALUES ($1, $2, $3)")
	findOneStmt = db.Prepare("SELECT id, size, rc FROM chunks WHERE id = $1")
	updateStmt = db.Prepare("UPDATE chunks SET rc = $1 WHERE id = $2")
	deleteStmt = db.Prepare("DELETE FROM chunks WHERE id = $1")
	statsStmt = db.Prepare("SELECT COUNT(*) AS count, TOTAL(size) as size FROM chunks")
}

func GetStats(ctx context.Context) (Stats, error) {
	var count uint64
	var totalSize float64
	if err := statsStmt.QueryRowContext(ctx).Scan(&count, &totalSize); err != nil {
		return Stats{}, fmt.Errorf("unable to query chunk stats: %v", err)
	}
	return Stats{
		Count:     uint64(count),
		TotalSize: uint64(totalSize),
	}, nil
}

func Create(ctx context.Context, data []byte) (string, error) {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	id, err := computeHash(data)
	if err != nil {
		return "", err
	}

	c, err := find(ctx, id)
	if err != nil {
		return "", err
	}
	if c == nil {
		if err := createNewChunk(ctx, id, data); err != nil {
			return "", err
		}
		return id, nil
	}

	if err := increaseReferenceCount(ctx, c); err != nil {
		return "", err
	}

	return id, nil
}

func Delete(ctx context.Context, id string) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	c, err := find(ctx, id)
	if err != nil {
		return err
	}

	if c == nil {
		return ErrNotFound
	}

	if c.References == 1 {
		if _, err := deleteStmt.ExecContext(ctx, c.ID); err != nil {
			return fmt.Errorf("unable to delete chunk %s: %v", c.ID, err)
		}
	}

	c.References -= 1

	if err := update(ctx, c); err != nil {
		return err
	}

	return nil
}

func createNewChunk(ctx context.Context, id string, data []byte) error {
	folder := id[:2]
	filename := id[2:]

	// create folder
	if err := os.Mkdir(path.Join(config.DataDir, "chunks", folder), 0700); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return err
		}
	}

	// write file
	if err := os.WriteFile(path.Join(config.DataDir, "chunks", folder, filename), data, 0700); err != nil {
		return err
	}

	// write table row
	if _, err := createStmt.ExecContext(ctx, id, len(data), 1); err != nil {
		return fmt.Errorf("unable to persist chunk: %v", err)
	}

	return nil
}

func increaseReferenceCount(ctx context.Context, c *Chunk) error {
	c.References += 1
	return update(ctx, c)
}

func find(ctx context.Context, id string) (*Chunk, error) {
	var chunk Chunk
	if err := findOneStmt.QueryRowContext(ctx, id).Scan(
		&chunk.ID,
		&chunk.Size,
		&chunk.References,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to check chunk existence: %v", err)
	}
	return &chunk, nil
}

func update(ctx context.Context, c *Chunk) error {
	if _, err := updateStmt.ExecContext(ctx, c.References, c.ID); err != nil {
		return fmt.Errorf("unable to update chunk record: %v", err)
	}
	return nil
}

func computeHash(data []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		return "", fmt.Errorf("unable to compute hash: %v", err)
	}
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

func Write(id string, w io.Writer) error {
	folder := id[:2]
	filename := id[2:]
	f, err := os.Open(path.Join(config.DataDir, "chunks", folder, filename))
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(w, f)
	return nil
}
