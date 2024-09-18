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

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
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
	ErrInvalidCredentials = jug.NewUnauthorizedError("invalid credentials")
	ErrNotFound           = jug.NewNotFoundError("api key not found")
	createStmt            *sql.Stmt
	listStmt              *sql.Stmt
	findStmt              *sql.Stmt
	getStmt               *sql.Stmt
	updateStmt            *sql.Stmt
	deleteStmt            *sql.Stmt
)

func Configure() {
	db.RunMigration("create_api_keys_table", `CREATE TABLE api_keys(
	id char(10) PRIMARY KEY,
	prefix char(10) NOT NULL,
	hash BLOB NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	created_by char(10) NOT NULL,
	expires_at TIMESTAMP NOT NULL
	)`)

	s := db.Prepare(
		"INSERT INTO api_keys (id, prefix, hash, description, created_at, created_by, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"SELECT * FROM api_keys ORDER BY created_at",
		"SELECT * FROM api_keys WHERE prefix = $1 LIMIT 1",
		"SELECT * FROM api_keys WHERE id = $1 LIMIT 1",
		"UPDATE api_keys SET expires_at = $1 WHERE id = $2",
		"DELETE FROM api_keys WHERE id = $1",
	)

	createStmt = s[0]
	listStmt = s[1]
	findStmt = s[2]
	getStmt = s[3]
	updateStmt = s[4]
	deleteStmt = s[5]
}

func Create(ctx context.Context, principal string, cmd CreateCommand) (*ApiKey, string, error) {
	now := time.Now()
	expiry := now.Add(cmd.TTL)
	prefix := domain.RandomId()
	suffix := domain.NewId(54)
	plainKey := prefix + suffix
	hash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("unable to hash API key: %v", err)
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
		return nil, "", fmt.Errorf("unable to save API key: %v", err)
	}

	return k, plainKey, nil
}

func List(ctx context.Context) ([]*ApiKey, error) {
	rows, err := listStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list api keys: %v", err)
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
			return nil, fmt.Errorf("unable to decode api key: %v", err)
		}
		keys = append(keys, &k)
	}
	return keys, nil
}

func Authenticate(ctx context.Context, key string) (*ApiKey, error) {
	if len(key) != 64 {
		return nil, ErrInvalidCredentials
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
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("unable to find api key: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(k.hash, []byte(key)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if k.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidCredentials
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
			return nil, ErrNotFound
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
