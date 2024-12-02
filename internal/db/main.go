// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/cfichtmueller/stor/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db                  *sql.DB
	findMigrationStmt   *sql.Stmt
	insertMigrationStmt *sql.Stmt
)

func Configure() {
	_db, err := openDb()
	if err != nil {
		log.Fatal(err)
	}

	db = _db

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id varchar(255) PRIMARY KEY, executed_at TIMESTAMP NOT NULL)"); err != nil {
		log.Fatalf("unable to create migration table: %v", err)
	}

	findMigrationStmt = Prepare("SELECT COUNT(*) as count FROM migrations WHERE id = $1")
	insertMigrationStmt = Prepare("INSERT INTO migrations (id, executed_at) VALUES ($1, $2)")

	runMigrations()
}

func Check() bool {
	fmt.Print("Checking database connection...")
	if _, err := openDb(); err != nil {
		fmt.Printf("\nERROR: unable to open database: %v\n", err)
		return false
	}
	fmt.Printf(" OK\n")
	return true
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

func QueryRow(query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func openDb() (*sql.DB, error) {
	dbUrl := "file:" + path.Join(config.DataDir, "db.s3db?mode=rwc")
	return sql.Open("sqlite3", dbUrl)
}
