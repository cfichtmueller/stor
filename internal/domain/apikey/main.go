// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apikey

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
	"golang.org/x/crypto/bcrypt"
)

type ApiKey struct {
	ID          string
	Prefix      string
	hash        []byte
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	ExpiresAt   time.Time
}

func (k *ApiKey) KeyMatches(key string) bool {
	if err := bcrypt.CompareHashAndPassword(k.hash, []byte(key)); err != nil {
		return false
	}
	return true
}

type CreateCommand struct {
	Description string
	TTL         time.Duration
}

var (
	createStmt *sql.Stmt
	listStmt   *sql.Stmt
	findStmt   *sql.Stmt
	getStmt    *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
)

func Configure() {
	createStmt = db.Prepare("INSERT INTO api_keys (id, prefix, hash, description, created_at, created_by, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	listStmt = db.Prepare("SELECT * FROM api_keys ORDER BY created_at")
	findStmt = db.Prepare("SELECT * FROM api_keys WHERE prefix = $1 LIMIT 1")
	getStmt = db.Prepare("SELECT * FROM api_keys WHERE id = $1 LIMIT 1")
	updateStmt = db.Prepare("UPDATE api_keys SET expires_at = $1 WHERE id = $2")
	deleteStmt = db.Prepare("DELETE FROM api_keys WHERE id = $1")
}

func Create(ctx context.Context, principal string, cmd CreateCommand) (*ApiKey, string, error) {
	now := domain.TimeNow()
	expiry := now.Add(cmd.TTL)
	prefix := domain.RandomId()
	suffix := domain.NewId(54)
	plainKey := prefix + suffix
	hash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("unable to hash API key: %w", err)
	}

	k := &ApiKey{
		ID:          domain.RandomId(),
		Prefix:      prefix,
		hash:        hash,
		Description: cmd.Description,
		CreatedAt:   now,
		CreatedBy:   principal,
		ExpiresAt:   expiry,
	}

	if _, err := createStmt.ExecContext(
		ctx,
		k.ID,
		k.Prefix,
		k.hash,
		k.Description,
		k.CreatedAt,
		k.CreatedBy,
		k.ExpiresAt,
	); err != nil {
		return nil, "", fmt.Errorf("unable to save API key: %w", err)
	}

	return k, plainKey, nil
}

func List(ctx context.Context) ([]*ApiKey, error) {
	rows, err := listStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list api keys: %w", err)
	}
	keys := make([]*ApiKey, 0)
	for rows.Next() {
		var k ApiKey
		if err := rows.Scan(
			&k.ID,
			&k.Prefix,
			&k.hash,
			&k.Description,
			&k.CreatedAt,
			&k.CreatedBy,
			&k.ExpiresAt,
		); err != nil {
			return nil, fmt.Errorf("unable to decode api key: %w", err)
		}
		keys = append(keys, &k)
	}
	return keys, nil
}

func Authenticate(ctx context.Context, key string) (*ApiKey, error) {
	if len(key) != 64 {
		return nil, ec.InvalidCredentials
	}
	prefix := key[:10]

	var k ApiKey
	if err := findStmt.QueryRowContext(ctx, prefix).Scan(
		&k.ID,
		&k.Prefix,
		&k.hash,
		&k.Description,
		&k.CreatedAt,
		&k.CreatedBy,
		&k.ExpiresAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.InvalidCredentials
		}
		return nil, fmt.Errorf("unable to find api key: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(k.hash, []byte(key)); err != nil {
		return nil, ec.InvalidCredentials
	}

	if k.ExpiresAt.Before(domain.TimeNow()) {
		return nil, ec.InvalidCredentials
	}

	return &k, nil
}

func Get(ctx context.Context, id string) (*ApiKey, error) {
	var k ApiKey
	if err := getStmt.QueryRowContext(ctx, id).Scan(
		&k.ID,
		&k.Prefix,
		&k.hash,
		&k.Description,
		&k.CreatedAt,
		&k.CreatedBy,
		&k.ExpiresAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.NoSuchApiKey
		}
		return nil, fmt.Errorf("unable to find api key: %v", err)
	}

	return &k, nil
}

func Update(ctx context.Context, key *ApiKey) error {
	if _, err := updateStmt.ExecContext(ctx, key.ExpiresAt, key.ID); err != nil {
		return fmt.Errorf("unable to update api key: %v", err)
	}
	return nil
}

func Delete(ctx context.Context, id string) error {
	if _, err := deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("unable to delete api key: %v", err)
	}
	return nil
}
