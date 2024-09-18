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

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/db"
)

type CreateCommand struct {
	Name string `json:"name"`
}

type Bucket struct {
	Name      string
	Objects   uint64
	Size      uint64
	CreatedAt time.Time
}

type Stats struct {
	Count        uint64
	TotalObjects uint64
}

func (b *Bucket) AddObject(size uint64) {
	b.Objects += 1
	b.Size += size
}

type Filter struct {
	Name string
}

var (
	ErrNotFound       = jug.NewNotFoundError("bucket not found")
	bucketNamePattern = regexp.MustCompile("^[a-z0-9](?:[a-z0-9.-]?[a-z0-9]+){2,}$")
	createStmt        *sql.Stmt
	findManyStmt      *sql.Stmt
	findOneStmt       *sql.Stmt
	updateStmt        *sql.Stmt
	statsStmt         *sql.Stmt
)

func Configure() {
	db.RunMigration("create_bucket_table", `CREATE TABLE buckets(
	name CHAR(64) PRIMARY KEY,
	objects INT NOT NULL,
	size INT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	created_by CHAR(32) NOT NULL
	)`)

	s := db.Prepare(
		"INSERT INTO buckets (name, objects, size, created_at, created_by) VALUES ($1, $2, $3, $4, $5)",
		"SELECT name, objects, size, created_at FROM buckets ORDER BY name ASC",
		"SELECT name, objects, size, created_at FROM buckets WHERE name = $1 LIMIT 1",
		"UPDATE buckets SET objects = $1, size = $2 WHERE name = $3",
		"SELECT COUNT(*) AS count, TOTAL(objects) AS objects from buckets",
	)
	createStmt = s[0]
	findManyStmt = s[1]
	findOneStmt = s[2]
	updateStmt = s[3]
	statsStmt = s[4]
}

func GetStats(ctx context.Context) (Stats, error) {
	var stats Stats
	if err := statsStmt.QueryRowContext(ctx).Scan(&stats.Count, &stats.TotalObjects); err != nil {
		return Stats{}, fmt.Errorf("unable to query bucket stats: %v", err)
	}
	return stats, nil
}

func Create(ctx context.Context, cmd CreateCommand) (*Bucket, error) {
	if !bucketNamePattern.MatchString(cmd.Name) || cmd.Name == "api" || cmd.Name == "css" || cmd.Name == "img" {
		return nil, jug.NewBadRequestError("invalid name")
	}
	b := &Bucket{
		Name:      cmd.Name,
		Objects:   0,
		Size:      0,
		CreatedAt: time.Now(),
	}
	if _, err := createStmt.ExecContext(ctx, b.Name, b.Objects, b.Size, b.CreatedAt, "system"); err != nil {
		return nil, fmt.Errorf("unable to create bucket record: %v", err)
	}
	return b, nil
}

func FindMany(ctx context.Context, filter *Filter) ([]*Bucket, error) {
	rows, err := findManyStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		buckets = append(buckets, &b)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buckets, nil
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
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unable to read db result: %v", err)
	}
	return &b, nil
}

func Save(ctx context.Context, b *Bucket) error {
	if _, err := updateStmt.ExecContext(ctx, b.Objects, b.Size, b.Name); err != nil {
		return fmt.Errorf("unable to save bucket: %v", err)
	}
	return nil
}
