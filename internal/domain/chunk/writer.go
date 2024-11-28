// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chunk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path"

	"github.com/cfichtmueller/stor/internal/domain"
)

type Writer struct {
	filename string
	hash     hash.Hash
	backing  io.WriteCloser
	size     int64
}

func NewWriter() (*Writer, error) {
	filename := path.Join(tempDir, domain.RandomId())
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &Writer{
		filename: filename,
		hash:     sha256.New(),
		backing:  file,
	}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if n, err := w.hash.Write(p); err != nil {
		return n, err
	}

	n, err = w.backing.Write(p)
	if err != nil {
		return n, err
	}
	w.size += int64(n)
	return n, nil
}

func (w *Writer) Close() error {
	return w.backing.Close()
}

// Commit commits the writer, resulting in a new chunk being created
func (w *Writer) Commit(ctx context.Context) (string, error) {
	hashBytes := w.hash.Sum(nil)
	id := hex.EncodeToString(hashBytes)

	writeMutex.Lock()
	defer writeMutex.Unlock()

	c, err := find(ctx, id)
	if err != nil {
		return "", err
	}

	if c == nil {
		filename, err := prepareChunkFile(id)
		if err != nil {
			return "", err
		}

		os.Rename(w.filename, filename)

		if err := createChunkTableRow(ctx, id, w.size); err != nil {
			return "", err
		}

		return id, nil
	}

	if err := IncreaseReferenceCount(ctx, c.ID); err != nil {
		return "", err
	}

	if err := os.Remove(w.filename); err != nil {
		return "", err
	}

	return id, nil
}

func (w *Writer) Size() int64 {
	return w.size
}
