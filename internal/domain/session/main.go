// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
)

type Session struct {
	ID         string
	User       string
	IpAddress  string
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

var (
	TTL         = time.Hour
	ErrNotFound = fmt.Errorf("session not found")
	createStmt  *sql.Stmt
	listStmt    *sql.Stmt
	getStmt     *sql.Stmt
	updateStmt  *sql.Stmt
	deleteStmt  *sql.Stmt
)

func Configure() {
	createStmt = db.Prepare("INSERT INTO sessions (id, user, ip_address, created_at, last_seen_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)")
	listStmt = db.Prepare("SELECT * FROM sessions WHERE user = $1 ORDER BY last_seen_at DESC")
	getStmt = db.Prepare("SELECT * FROM sessions where id = $1 LIMIT 1")
	updateStmt = db.Prepare("UPDATE sessions SET last_seen_at = $1, expires_at = $2 WHERE id = $3")
	deleteStmt = db.Prepare("DELETE FROM sessions WHERE id = $1")
}

func Create(ctx context.Context, user, ipAddress string) (*Session, error) {
	now := domain.TimeNow()
	s := &Session{
		ID:         domain.NewId(64),
		User:       user,
		IpAddress:  ipAddress,
		CreatedAt:  now,
		LastSeenAt: now,
		ExpiresAt:  now.Add(TTL),
	}

	if _, err := createStmt.ExecContext(ctx, s.ID, s.User, s.IpAddress, s.CreatedAt, s.LastSeenAt, s.ExpiresAt); err != nil {
		return nil, fmt.Errorf("unable to create session record: %w", err)
	}

	return s, nil
}

func List(ctx context.Context, user string) ([]*Session, error) {
	rows, err := listStmt.QueryContext(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("unable to query sessions: %w", err)
	}
	res := make([]*Session, 0)
	for rows.Next() {
		s, err := decode(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func Get(ctx context.Context, id string) (*Session, error) {
	row := getStmt.QueryRowContext(ctx, id)
	s, err := decode(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func MarkSeen(ctx context.Context, id string) error {
	now := domain.TimeNow()
	if _, err := updateStmt.ExecContext(ctx, now, now.Add(TTL), id); err != nil {
		return fmt.Errorf("unable to update session record: %w", err)
	}
	return nil
}

func Delete(ctx context.Context, id string) error {
	if _, err := deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("unable to delete session: %w", err)
	}
	return nil
}

func decode(row Scanner) (*Session, error) {
	s := Session{}
	if err := row.Scan(
		&s.ID,
		&s.User,
		&s.IpAddress,
		&s.CreatedAt,
		&s.LastSeenAt,
		&s.ExpiresAt,
	); err != nil {
		return nil, fmt.Errorf("unable to decode session: %w", err)
	}
	return &s, nil
}

type Scanner interface {
	Scan(dest ...any) error
}
