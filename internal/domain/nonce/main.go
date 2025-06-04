// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nonce

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
)

type CreateCommand struct {
	TTL time.Duration
}

type Nonce struct {
	ID        string
	Bucket    string
	Key       string
	ExpiresAt time.Time
}

var (
	ErrNotFound       = fmt.Errorf("nonce not found")
	createStmt        *sql.Stmt
	findOneStmt       *sql.Stmt
	deleteStmt        *sql.Stmt
	deleteExpiredStmt *sql.Stmt
)

func Configure() {
	createStmt = db.Prepare("INSERT INTO nonces (id, bucket, key, expires_at) VALUES ($1, $2, $3, $4)")
	findOneStmt = db.Prepare("SELECT id, bucket, key, expires_at FROM nonces WHERE id = $1 LIMIT 1")
	deleteStmt = db.Prepare("DELETE FROM nonces WHERE id = $1")
	deleteExpiredStmt = db.Prepare("DELETE FROM nonces WHERE expires_at < $1")

	go worker()
}

func Create(ctx context.Context, bucket, key string, cmd CreateCommand) (*Nonce, error) {
	nonce := &Nonce{
		ID:        domain.NewId(64),
		Bucket:    bucket,
		Key:       key,
		ExpiresAt: domain.TimeNow().Add(cmd.TTL),
	}
	if _, err := createStmt.ExecContext(ctx, nonce.ID, nonce.Bucket, nonce.Key, nonce.ExpiresAt); err != nil {
		return nil, fmt.Errorf("unable to create nonce record: %w", err)
	}
	return nonce, nil
}

func Get(ctx context.Context, id string) (*Nonce, error) {
	var nonce Nonce
	if err := findOneStmt.QueryRowContext(ctx, id).Scan(&nonce.ID, &nonce.Bucket, &nonce.Key, &nonce.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unable to find nonce: %w", err)
	}
	if nonce.ExpiresAt.Before(time.Now()) {
		return nil, ErrNotFound
	}
	return &nonce, nil
}

func GetAndInvalidate(ctx context.Context, id string) (*Nonce, error) {
	n, err := Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := deleteStmt.ExecContext(ctx, id); err != nil {
		return nil, fmt.Errorf("unable to delete nonce: %w", err)
	}

	return n, nil
}

func worker() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		purgeNonces()
	}
}

func purgeNonces() {
	if _, err := deleteExpiredStmt.Exec(domain.TimeNow()); err != nil {
		slog.Error("unable to purge expired nonces", "error", err)
	}
}
