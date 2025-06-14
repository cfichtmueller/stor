// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bucket

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
)

type CreateCommand struct {
	Name string `json:"name"`
}

type Bucket struct {
	Name      string
	Objects   int64
	Size      int64
	CreatedAt time.Time
}

type Stats struct {
	Count        int
	TotalObjects int
}

type Filter struct {
	Name string
}

var (
	bucketNamePattern = regexp.MustCompile("^[a-z0-9](?:[a-z0-9.-]?[a-z0-9]+){2,}$")
	createStmt        *sql.Stmt
	findManyStmt      *sql.Stmt
	findOneStmt       *sql.Stmt
	updateStmt        *sql.Stmt
	statsStmt         *sql.Stmt
	listStmt          *sql.Stmt
	countStmt         *sql.Stmt
	deleteStmt        *sql.Stmt
)

func Configure() {
	createStmt = db.Prepare("INSERT INTO buckets (name, objects, size, created_at, created_by) VALUES ($1, $2, $3, $4, $5)")
	findManyStmt = db.Prepare("SELECT name, objects, size, created_at FROM buckets ORDER BY name ASC")
	findOneStmt = db.Prepare("SELECT name, objects, size, created_at FROM buckets WHERE name = $1 LIMIT 1")
	updateStmt = db.Prepare("UPDATE buckets SET objects = $1, size = $2 WHERE name = $3")
	statsStmt = db.Prepare("SELECT COUNT(*) AS count, TOTAL(objects) AS objects from buckets")
	listStmt = db.Prepare("SELECT name, objects, size, created_at FROM buckets WHERE name > $1 ORDER BY name LIMIT $2")
	countStmt = db.Prepare("SELECT COUNT(*) FROM buckets WHERE name > $1")
	deleteStmt = db.Prepare("DELETE FROM buckets WHERE name = $1")
}

func ValidateName(name string) error {
	v := srv.RequireNotEmpty("name", name, nil)
	v = srv.RequireMaxLength("name", 63, name, v)
	v = srv.Require("name", srv.ValidationCodeInvalid, "name must be a valid bucket name", name != "api" && name != "css" && name != "img", v)
	if v == nil {
		v = srv.RequireRegex("name", name, bucketNamePattern, nil)
	}
	return srv.Validate(v)
}

func GetStats(ctx context.Context) (Stats, error) {
	var stats Stats
	if err := statsStmt.QueryRowContext(ctx).Scan(&stats.Count, &stats.TotalObjects); err != nil {
		return Stats{}, fmt.Errorf("unable to query bucket stats: %w", err)
	}
	return stats, nil
}

func Create(ctx context.Context, cmd CreateCommand) (*Bucket, error) {
	b := &Bucket{
		Name:      cmd.Name,
		Objects:   0,
		Size:      0,
		CreatedAt: domain.TimeNow(),
	}
	if _, err := createStmt.ExecContext(ctx, b.Name, b.Objects, b.Size, b.CreatedAt, "system"); err != nil {
		return nil, fmt.Errorf("unable to create bucket record: %w", err)
	}
	return b, nil
}

func List(ctx context.Context, startAfter string, maxBuckets int) ([]*Bucket, error) {
	return decodeRows(listStmt.QueryContext(ctx, startAfter, maxBuckets))
}

func FindMany(ctx context.Context, filter *Filter) ([]*Bucket, error) {
	return decodeRows(findManyStmt.QueryContext(ctx))
}

func FindOne(ctx context.Context, name string) (*Bucket, error) {
	var b Bucket
	if err := findOneStmt.QueryRowContext(ctx, name).
		Scan(
			&b.Name,
			&b.Objects,
			&b.Size,
			&b.CreatedAt,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.NoSuchBucket
		}
		return nil, fmt.Errorf("unable to read db result: %w", err)
	}
	return &b, nil
}

func Count(ctx context.Context, startAfter string) (int, error) {
	var count int
	if err := countStmt.QueryRowContext(ctx, startAfter).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to query bucket count: %w", err)
	}
	return count, nil
}

func Save(ctx context.Context, b *Bucket) error {
	if _, err := updateStmt.ExecContext(ctx, b.Objects, b.Size, b.Name); err != nil {
		return fmt.Errorf("unable to save bucket: %w", err)
	}
	return nil
}

func Delete(ctx context.Context, name string) error {
	if _, err := deleteStmt.ExecContext(ctx, name); err != nil {
		return fmt.Errorf("unable to delete bucket: %w", err)
	}
	return nil
}

func decodeRows(rows *sql.Rows, err error) ([]*Bucket, error) {
	if err != nil {
		return nil, fmt.Errorf("unable to query buckets: %w", err)
	}
	buckets := make([]*Bucket, 0)
	for rows.Next() {
		var b Bucket
		if err := rows.Scan(
			&b.Name,
			&b.Objects,
			&b.Size,
			&b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		buckets = append(buckets, &b)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("unable to close rows: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}
	return buckets, nil
}
