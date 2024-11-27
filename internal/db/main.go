// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"database/sql"
	"log"
	"path"
	"time"

	"github.com/cfichtmueller/stor/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db                  *sql.DB
	findMigrationStmt   *sql.Stmt
	insertMigrationStmt *sql.Stmt
)

func Configure() {
	dbUrl := "file:" + path.Join(config.DataDir, "db.s3db?mode=rwc")
	_db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	db = _db

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id varchar(255) PRIMARY KEY, executed_at TIMESTAMP NOT NULL)"); err != nil {
		log.Fatalf("unable to create migration table: %v", err)
	}

	findMigrationStmt = Prepare("SELECT COUNT(*) as count FROM migrations WHERE id = $1")
	insertMigrationStmt = Prepare("INSERT INTO migrations (id, executed_at) VALUES ($1, $2)")
}

func Exec(query string, args ...any) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Prepare(statement string) *sql.Stmt {
	s, err := db.Prepare(statement)
	if err != nil {
		log.Fatalf("unable to prepare statement '%s': %v", statement, err)
	}
	return s
}

func PrepareOne(query string) (*sql.Stmt, error) {
	return db.Prepare(query)
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryRow(query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

func RunMigration(id, statement string) {
	RunMigrationF(id, func() error {
		ctx := context.Background()
		_, err := db.ExecContext(ctx, statement)
		return err
	})
}

func RunMigrationF(id string, f func() error) {
	ctx := context.Background()
	var count int
	if err := findMigrationStmt.QueryRowContext(ctx, id).Scan(&count); err != nil {
		log.Fatalf("unable to query migration: %v", err)
	}
	if count > 0 {
		return
	}
	if err := f(); err != nil {
		log.Fatalf("unable to run migration %s: %v", id, err)
	}
	if _, err := insertMigrationStmt.Exec(id, time.Now()); err != nil {
		log.Fatalf("unable to persist migration status: %v", err)
	}
}
