package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cfichtmueller/stor/internal/domain"
)

func runMigrations() {
	m("create_api_keys_table", `CREATE TABLE api_keys(
		id char(10) PRIMARY KEY,
		prefix char(10) NOT NULL,
		hash BLOB NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		created_by char(10) NOT NULL,
		expires_at TIMESTAMP NOT NULL
	)`)
	m("create_bucket_table", `CREATE TABLE buckets(
		name CHAR(64) PRIMARY KEY,
		objects INT NOT NULL,
		size INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		created_by CHAR(32) NOT NULL
	)`)
	m("create_chunks_table", `CREATE TABLE chunks(
		id CHAR(64) PRIMARY KEY,
		size INT NOT NULL,
		rc INT NOT NULL
	)`)
	m("create_object_table", `CREATE TABLE objects(
		id CHAR(32) PRIMARY KEY,
		bucket CHAR(64) NOT NULL,
		key TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INT NOT NULL,
		created_at DATETIME NOT NULL
	)`)
	m("create_object_chunks_table", `CREATE TABLE object_chunks(
		object CHAR(32),
		chunk CHAR(64),
		seq INT,
		PRIMARY KEY (object, chunk, seq)
	)`)
	m("create_sessions_table", `CREATE TABLE sessions(
		id CHAR(64) PRIMARY KEY,
		user CHAR(10) NOT NULL,
		ip_address CHAR(40) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		last_seen_at TIMESTAMP NOT NULL,
		expires_at TIMESTAMP NOT NULL
	)`)
	m("create_users_table", `CREATE TABLE users(
		id CHAR(10) PRIMARY KEY,
		email TEXT NOT NULL,
		enabled BOOLEAN NOT NULL,
		password_hash BLOB NOT NULL,
		created_at TIMESTAMP NOT NULL,
		last_seen_at TIMESTAMP
	)`)

	m("add_object_deleted_flag", `ALTER TABLE objects ADD COLUMN is_deleted INTEGER`)
	m("add_object_key_index", `CREATE INDEX idx_objects_key_bucket_deleted ON objects (key, bucket, is_deleted)`)
	m("add_objectchunk_key_index", `CREATE INDEX idx_objectchunk_key ON object_chunks (object)`)

	// object etag setup
	m("add_object_etag_1", `ALTER TABLE objects ADD COLUMN etag CHAR(64)`)
	mf("add_object_etags", func() error {
		find := Prepare("SELECT id, key FROM objects WHERE etag IS NULL AND key > $1 ORDER BY key LIMIT 10000")
		update := Prepare("UPDATE objects SET etag = $1 WHERE id = $2")
		start := ""
		for {
			rows, err := find.Query(start)
			if err != nil {
				return err
			}
			ids := make([]string, 0)
			for rows.Next() {
				var id, key string
				if err := rows.Scan(&id, &key); err != nil {
					return err
				}
				ids = append(ids, id)
				start = key
			}
			if len(ids) == 0 {
				break
			}
			for _, id := range ids {
				if _, err := update.Exec(domain.NewEtag(), id); err != nil {
					return err
				}
			}
		}
		return nil
	})

	// archives setup
	m("create_archive_table", `CREATE TABLE archives (
		id CHAR(32) PRIMARY KEY,
		bucket CHAR(64) NOT NULL,
		key TEXT NOT NULL,
		type CHAR(6) NOT NULL,
		state CHAR(64) NOT NULL, 
		is_deleted INTEGER NOT NULL
	)`)
	m("create_archive_entries_table", `CREATE TABLE archive_entries (
		id CHAR(32) PRIMARY KEY,
		archive CHAR(32) NOT NULL,
		key TEXT NOT NULL,
		name TEXT NOT NULL
	)`)
	m("create_archive_entries_index", `CREATE INDEX idx_archive_entries ON archive_entries (archive, name)`)

	// nonces setup
	m("create_nonces_table", `CREATE TABLE nonces(
		id CHAR(64) PRIMARY KEY,
		bucket CHAR(64) NOT NULL,
		key TEXT NOT NULL,
		expires_at DATETIME NOT NULL
	)`)

	// object versions setup
	mf("add_object_versions", func() error {
		if _, err := db.Exec(`ALTER TABLE objects ADD COLUMN current CHAR(32)`); err != nil {
			return err
		}
		if _, err := db.Exec(`CREATE TABLE object_versions(
		id CHAR(32) PRIMARY KEY,
		object CHAR(32) NOT NULL,
		content_type TEXT NOT NULL,
		size INT NOT NULL,
		created_at DATETIME NOT NULL,
		etag CHAR(64) NOT NULL,
		is_deleted INT NOT NULL
	)`); err != nil {
			return err
		}

		type mo struct {
			id          string
			contentType string
			size        int64
			createdAt   time.Time
			etag        string
		}

		offset := 0
		for {
			rows, err := db.Query("SELECT id, content_type, size, created_at, etag FROM objects ORDER BY id LIMIT 1000 OFFSET ?", offset)
			if err != nil {
				return err
			}

			mos := make([]*mo, 0)

			for rows.Next() {
				e := mo{}
				if err := rows.Scan(
					&e.id,
					&e.contentType,
					&e.size,
					&e.createdAt,
					&e.etag,
				); err != nil {
					return fmt.Errorf("unable to decode object row: %v", err)
				}
				mos = append(mos, &e)
			}

			if len(mos) == 0 {
				break
			}
			offset += len(mos)

			for _, e := range mos {
				versionId := domain.RandomId()
				if _, err := db.Exec("INSERT INTO object_versions (id, object, content_type, size, created_at, etag, is_deleted) VALUES (?, ?, ?, ?, ?, ?, 0)",
					versionId,
					e.id,
					e.contentType,
					e.size,
					e.createdAt,
					e.etag,
				); err != nil {
					return fmt.Errorf("unable to create object version for object %s: %v", e.id, err)
				}
				if _, err := db.Exec("UPDATE objects SET current = ? WHERE id = ?", versionId, e.id); err != nil {
					return fmt.Errorf("unable to set object version pointer for object %s: %v", e.id, err)
				}
				if _, err := db.Exec("UPDATE object_chunks SET object = ? WHERE object = ?", versionId, e.id); err != nil {
					return fmt.Errorf("unable to update object chunks object pointers for object %s: %v", e.id, err)
				}
			}
		}
		return nil
	})

	// api key created by setup
	mf("20241127_initialize_api_key_created_by", func() error {
		if _, err := db.Exec(`CREATE TABLE api_keys_dg_tmp(
			id CHAR(10) PRIMARY KEY,
			prefix      CHAR(10)  NOT NULL,
			hash        BLOB      NOT NULL,
			description TEXT      NOT NULL,
			created_at  TIMESTAMP NOT NULL,
			created_by  CHAR(32)  NOT NULL,
			expires_at  TIMESTAMP NOT NULL
		)`); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT INTO api_keys_dg_tmp(id, prefix, hash, description, created_at, created_by, expires_at)
			SELECT id, prefix, hash, description, created_at, created_by, expires_at FROM api_keys;`); err != nil {
			return err
		}
		if _, err := db.Exec("DROP TABLE api_keys"); err != nil {
			return err
		}
		if _, err := db.Exec("ALTER TABLE api_keys_dg_tmp RENAME TO api_keys;"); err != nil {
			return err
		}
		var id string
		if err := db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		if _, err := db.Exec("UPDATE api_keys SET created_by = ?", "user:"+id); err != nil {
			return fmt.Errorf("unable to update api keys: %v", err)
		}
		return nil
	})
}

func m(id, statement string) {
	mf(id, func() error {
		ctx := context.Background()
		_, err := db.ExecContext(ctx, statement)
		return err
	})
}

func mf(id string, f func() error) {
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
