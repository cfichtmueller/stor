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
	"path/filepath"
	"strings"
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
	ErrNotFound                = fmt.Errorf("chunk not found")
	tempDir                    string
	createStmt                 *sql.Stmt
	findOneStmt                *sql.Stmt
	updateStmt                 *sql.Stmt
	decreaseReferenceCountStmt *sql.Stmt
	increaseReferenceCountStmt *sql.Stmt
	deleteStmt                 *sql.Stmt
	statsStmt                  *sql.Stmt
	writeMutex                 sync.Mutex
)

func Configure() {
	if err := os.Mkdir(path.Join(config.DataDir, "chunks"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("unable to create chunk directory: %v", err)
		}
	}

	tempDir = path.Join(config.DataDir, "chunk_tmp")
	if err := os.Mkdir(tempDir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("unable to create chunk temp directory: %v", err)
	}

	createStmt = db.Prepare("INSERT INTO chunks (id, size, rc) VALUES ($1, $2, $3)")
	findOneStmt = db.Prepare("SELECT id, size, rc FROM chunks WHERE id = $1")
	updateStmt = db.Prepare("UPDATE chunks SET rc = $1 WHERE id = $2")
	decreaseReferenceCountStmt = db.Prepare("UPDATE chunks SET rc = rc - 1 WHERE id = ?")
	increaseReferenceCountStmt = db.Prepare("UPDATE chunks SET rc = rc + 1 WHERE id = ?")
	deleteStmt = db.Prepare("DELETE FROM chunks WHERE id = $1")
	statsStmt = db.Prepare("SELECT COUNT(*) AS count, TOTAL(size) as size FROM chunks")
}

func Check() bool {
	Configure()
	fmt.Println("Checking chunks...")
	var count int64
	var size float64
	if err := statsStmt.QueryRow().Scan(&count, &size); err != nil {
		fmt.Printf("ERROR: unable to get chunk statistics: %v\n", err)
		return false
	}
	fmt.Printf("Found %d chunks in database. Total size: %d\n", count, int64(size))
	index := make(map[string]map[string]struct{})
	offset := 0
	for {
		r, err := db.Query("SELECT id FROM chunks ORDER BY id LIMIT 1000 OFFSET ?", offset)
		if err != nil {
			fmt.Printf("ERROR: unable to scan chunks: %v\n", err)
			return false
		}
		read := false
		for r.Next() {
			read = true
			offset += 1
			var id string
			if err := r.Scan(&id); err != nil {
				fmt.Printf("ERROR: unable to decode chunk row: %v\n", err)
				return false
			}
			prefixes, ok := index[id[:2]]
			if !ok {
				v := make(map[string]struct{})
				v[id[2:]] = struct{}{}
				index[id[:2]] = v
			} else {
				prefixes[id[2:]] = struct{}{}
			}
		}
		if !read {
			break
		}
	}
	//TODO: also check for missing chunk files
	fmt.Printf("Scanning chunk files...")
	dangling := make([]string, 0)
	chunksFolder := path.Join(config.DataDir, "chunks")
	separator := string(os.PathSeparator)
	base := chunksFolder + separator
	nok := false
	log := func(format string, args ...any) {
		if !nok {
			fmt.Print("\n")
			nok = true
		}
		fmt.Printf(format, args...)
	}
	if err := filepath.Walk(chunksFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log("ERROR\n  unable to access path %s: %v\n", path, err)
			return err
		}
		if path == chunksFolder {
			return nil
		}
		trimmedPath := strings.TrimPrefix(path, base)
		parts := strings.Split(trimmedPath, separator)

		prefix := parts[0]
		if len(parts) == 1 { // chunk folder
			_, ok := index[prefix]
			if !ok {
				log("WARNING: found dangling chunk folder %s\n", trimmedPath)
				dangling = append(dangling, path)
			}
		} else { // chunk
			if _, ok := index[prefix]; !ok {
				log("WARNING: found dangling chunk %s\n", trimmedPath)
				dangling = append(dangling, path)
				return nil
			}
			if _, ok := index[prefix][parts[1]]; !ok {
				log("WARNING: found dangling chunk %s\n", trimmedPath)
				dangling = append(dangling, path)
				return nil
			}
		}
		return nil
	}); err != nil {
		return false
	}
	if nok {
		fmt.Print("Done scanning chunk files\n")
	} else {
		fmt.Print("OK\n")
	}

	return true
}

func GetStats(ctx context.Context) (Stats, error) {
	var count uint64
	var totalSize float64
	if err := statsStmt.QueryRowContext(ctx).Scan(&count, &totalSize); err != nil {
		return Stats{}, fmt.Errorf("unable to query chunk stats: %w", err)
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

	if err := IncreaseReferenceCount(ctx, c.ID); err != nil {
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

	if c.References > 1 {
		return DecreaseReferenceCount(ctx, id)
	}

	if _, err := deleteStmt.ExecContext(ctx, c.ID); err != nil {
		return fmt.Errorf("unable to delete chunk %s: %w", c.ID, err)
	}
	folder := id[:2]
	filename := id[2:]
	if err := os.Remove(path.Join(config.DataDir, "chunks", folder, filename)); err != nil {
		return fmt.Errorf("unable to delete chunk file: %w", err)
	}
	return nil
}

func IncreaseReferenceCount(ctx context.Context, chunkId string) error {
	if _, err := increaseReferenceCountStmt.ExecContext(ctx, chunkId); err != nil {
		return fmt.Errorf("unable to increase reference count for chunk %s: %w", chunkId, err)
	}
	return nil
}

func DecreaseReferenceCount(ctx context.Context, chunkId string) error {
	if _, err := decreaseReferenceCountStmt.ExecContext(ctx, chunkId); err != nil {
		return fmt.Errorf("unable to decrease reference count for chunk %s: %w", chunkId, err)
	}
	return nil
}

func createNewChunk(ctx context.Context, id string, data []byte) error {
	filename, err := prepareChunkFile(id)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0700); err != nil {
		return err
	}

	if err := createChunkTableRow(ctx, id, int64(len(data))); err != nil {
		return err
	}

	return nil
}

// prepareChunkFile prepares the chunk folder and filename. The method returns the filename
func prepareChunkFile(id string) (string, error) {
	folder := id[:2]
	filename := id[2:]

	// create folder
	if err := os.Mkdir(path.Join(config.DataDir, "chunks", folder), 0700); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return "", err
		}
	}
	return path.Join(config.DataDir, "chunks", folder, filename), nil
}

func createChunkTableRow(ctx context.Context, id string, size int64) error {
	if _, err := createStmt.ExecContext(ctx, id, size, 1); err != nil {
		return fmt.Errorf("unable to persist chunk: %w", err)
	}
	return nil
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
		return nil, fmt.Errorf("unable to check chunk existence: %w", err)
	}
	return &chunk, nil
}

func update(ctx context.Context, c *Chunk) error {
	if _, err := updateStmt.ExecContext(ctx, c.References, c.ID); err != nil {
		return fmt.Errorf("unable to update chunk record: %w", err)
	}
	return nil
}

func computeHash(data []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		return "", fmt.Errorf("unable to compute hash: %w", err)
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
