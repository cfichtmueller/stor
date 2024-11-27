// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package object

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/ec"
)

type CreateCommand struct {
	Key         string
	ContentType string
	Data        []byte
	Size        int64
}

type UpdateCommand struct {
	ContentType string
	Data        []byte
	Size        int64
}

type Object struct {
	ID             string
	Bucket         string
	Key            string
	ETag           string
	ContentType    string
	Size           int64
	CreatedAt      time.Time
	Deleted        bool
	CurrentVersion string
}

const (
	objectsTable        = "objects"
	objectVersionsTable = "object_versions"
	objectChunksTable   = "object_chunks"
	chunksTable         = "chunks"
)

var (
	purgeMutex   sync.Mutex
	maxPurgeTime = 400 * time.Millisecond
	purgeFlag    = true

	objectFields = "id, bucket, key, etag, content_type, size, created_at, is_deleted, current"

	createStmt               *sql.Stmt
	listStmt                 *sql.Stmt
	findOneStmt              *sql.Stmt
	existsStmt               *sql.Stmt
	updateObjectMetadataStmt *sql.Stmt
	// Marks an object as deleted. Input: object id
	markObjectDeletedStmt *sql.Stmt
	// Finds all deleted objects. Input: none
	findDeletedObjectsStmt *sql.Stmt
	// Deletes an object. Input: object id
	deleteObjectStmt     *sql.Stmt
	addObjectChunkStmt   *sql.Stmt
	findObjectChunksStmt *sql.Stmt
	//TODO: rename object col to version
	// Deletes all object chunks ob an object. Input: object id
	deleteObjectChunksStmt  *sql.Stmt
	countStmt               *sql.Stmt
	statsStmt               *sql.Stmt
	createObjectVersionStmt *sql.Stmt
	// Marks all object versions of an object as deleted. Input: object id
	markObjectVersionsDeletedStmt *sql.Stmt
	// Marks an object version as deleted. Input: object version id
	markObjectVersionDeletedStmt *sql.Stmt
	// Finds all deleted objects. Input: none
	findDeletedObjectVersionsStmt *sql.Stmt
	// Deletes an object version. Input: object version id
	deleteObjectVersionStmt *sql.Stmt
)

func Configure() {
	if err := config.Mkdir("chunks"); err != nil {
		log.Fatalf("unable to create chunk directory: %v", err)
	}

	createStmt = db.Prepare("INSERT INTO objects (" + objectFields + ") VALUES ($1, $2, $3, $4, $5, $6, $7, false, $8)")
	listStmt = db.Prepare("SELECT " + objectFields + " FROM objects WHERE bucket = $1 AND key > $2 AND is_deleted = $3 ORDER BY key LIMIT $4")
	findOneStmt = db.Prepare("SELECT " + objectFields + " FROM objects WHERE bucket = $1 AND key = $2 AND is_deleted = $3 LIMIT 1")
	existsStmt = db.Prepare("SELECT COUNT(*) as count FROM objects WHERE bucket = $1 AND key = $2 AND is_deleted = $3")
	markObjectDeletedStmt = db.Prepare("UPDATE objects SET is_deleted = 1 WHERE id = ?")
	deleteObjectStmt = db.Prepare("DELETE FROM objects WHERE id = ?")
	findDeletedObjectsStmt = db.Prepare("SELECT id FROM objects WHERE is_deleted = true LIMIT 1000")
	addObjectChunkStmt = db.Prepare("INSERT INTO object_chunks (object, chunk, seq) VALUES ($1, $2, $3)")
	findObjectChunksStmt = db.Prepare("SELECT chunk FROM object_chunks WHERE object = $1 ORDER BY seq")
	deleteObjectChunksStmt = db.Prepare("DELETE FROM object_chunks WHERE object = $1")
	countStmt = db.Prepare("SELECT COUNT(*) FROM objects WHERE bucket = $1 AND key > $2 AND is_deleted  = $3")
	statsStmt = db.Prepare("SELECT COUNT(*), TOTAL(size) FROM objects WHERE bucket = $1 AND is_deleted = $2")
	createObjectVersionStmt = db.Prepare("INSERT INTO object_versions (id, object, content_type, size, created_at, etag, is_deleted) VALUES (?, ?, ?, ?, ?, ?, 0)")
	updateObjectMetadataStmt = db.Prepare("UPDATE objects SET content_type = ?, size = ?, etag = ?, current = ? WHERE id = ?")
	markObjectVersionsDeletedStmt = db.Prepare("UPDATE object_versions SET is_deleted = 1 WHERE object = ?")
	markObjectVersionDeletedStmt = db.Prepare("UPDATE object_versions SET is_deleted = 1 WHERE id = ?")
	findDeletedObjectVersionsStmt = db.Prepare("SELECT id FROM object_versions WHERE is_deleted = true LIMIT 1000")
	deleteObjectVersionStmt = db.Prepare("DELETE FROM object_versions WHERE id = ?")

	go worker()
}

func List(ctx context.Context, bucketName, startAfter string, limit int) ([]*Object, error) {
	return decodeRows(listStmt.QueryContext(ctx, bucketName, startAfter, false, limit))
}

func Count(ctx context.Context, bucketName, startAfter string) (int, error) {
	var count int
	if err := countStmt.QueryRowContext(ctx, bucketName, startAfter, false).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to count objects: %v", err)
	}
	return count, nil
}

type Stats struct {
	ObjectCount int64
	TotalSize   int64
}

func StatsForBucket(ctx context.Context, bucketName string) (*Stats, error) {
	var objects int64
	var size float64
	if err := statsStmt.QueryRowContext(ctx, bucketName, false).Scan(&objects, &size); err != nil {
		return nil, fmt.Errorf("unable to query object stats: %v", err)
	}

	return &Stats{
		ObjectCount: objects,
		TotalSize:   int64(size),
	}, nil
}

func decodeRows(rows *sql.Rows, err error) ([]*Object, error) {
	if err != nil {
		return nil, fmt.Errorf("unable to find object records: %v", err)
	}
	objects := make([]*Object, 0)
	for rows.Next() {
		var o Object
		if err := rows.Scan(
			&o.ID,
			&o.Bucket,
			&o.Key,
			&o.ETag,
			&o.ContentType,
			&o.Size,
			&o.CreatedAt,
			&o.Deleted,
			&o.CurrentVersion,
		); err != nil {
			return nil, fmt.Errorf("unable to decode object record: %v", err)
		}
		objects = append(objects, &o)
	}
	return objects, nil
}

// FindOne finds an object. Returns ec.NoSuchKey if the object cannot be found
func FindOne(ctx context.Context, bucketName, key string, deleted bool) (*Object, error) {
	var o Object
	if err := findOneStmt.QueryRowContext(ctx, bucketName, key, deleted).Scan(
		&o.ID,
		&o.Bucket,
		&o.Key,
		&o.ETag,
		&o.ContentType,
		&o.Size,
		&o.CreatedAt,
		&o.Deleted,
		&o.CurrentVersion,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ec.NoSuchKey
		}
		return nil, fmt.Errorf("unable to find object record: %v", err)
	}
	return &o, nil
}

func FindMany(ctx context.Context, bucketName string, keys []string, deleted bool) ([]*Object, error) {
	query := strings.Builder{}
	query.WriteString("SELECT " + objectFields + " FROM objects WHERE bucket = $1 AND is_deleted = $2 AND key IN (")
	first := true
	for i := range keys {
		if first {
			first = false
		} else {
			query.WriteString(", ")
		}
		query.WriteString("$" + strconv.Itoa(i+3))
	}
	query.WriteString(")")
	stmt, err := db.PrepareOne(query.String())
	if err != nil {
		return nil, err
	}
	args := make([]any, 0, len(keys)+3)
	args = append(args, bucketName)
	args = append(args, deleted)
	for _, k := range keys {
		args = append(args, k)
	}
	return decodeRows(stmt.QueryContext(ctx, args...))
}

func Exists(ctx context.Context, bucketName, key string) (bool, error) {
	var count int
	if err := existsStmt.QueryRowContext(ctx, bucketName, key, false).Scan(&count); err != nil {
		return false, fmt.Errorf("unable to count objects: %v", err)
	}
	return count > 0, nil
}

func Create(ctx context.Context, bucketId string, cmd CreateCommand) (*Object, error) {
	chunkId, err := chunk.Create(ctx, cmd.Data)
	if err != nil {
		return nil, err
	}

	return CreateWithChunk(ctx, bucketId, chunkId, cmd)
}

func CreateWithChunk(ctx context.Context, bucketId, chunkId string, cmd CreateCommand) (*Object, error) {
	size := cmd.Size
	if size == 0 && cmd.Data != nil {
		size = int64(len(cmd.Data))
	}
	o := Object{
		ID:             domain.RandomId(),
		Bucket:         bucketId,
		Key:            cmd.Key,
		ETag:           domain.NewEtag(),
		ContentType:    cmd.ContentType,
		Size:           size,
		CreatedAt:      domain.TimeNow(),
		CurrentVersion: domain.RandomId(),
	}

	if _, err := createObjectVersionStmt.ExecContext(ctx, o.CurrentVersion, o.ID, o.ContentType, o.Size, o.CreatedAt, o.ETag); err != nil {
		return nil, fmt.Errorf("unable to create object version: %v", err)
	}
	if _, err := addObjectChunkStmt.ExecContext(ctx, o.CurrentVersion, chunkId, 1); err != nil {
		return nil, fmt.Errorf("unable to persist object chunk record: %v", err)
	}
	if _, err := createStmt.ExecContext(ctx, o.ID, o.Bucket, o.Key, o.ETag, o.ContentType, o.Size, o.CreatedAt, o.CurrentVersion); err != nil {
		return nil, fmt.Errorf("unable to persist object record: %v", err)
	}

	return &o, nil
}

var updateLock sync.Mutex

func Update(ctx context.Context, o *Object, cmd UpdateCommand) (*Object, error) {
	size := cmd.Size
	if size == 0 && cmd.Data != nil {
		size = int64(len(cmd.Data))
	}
	chunkId, err := chunk.Create(ctx, cmd.Data)
	if err != nil {
		return nil, err
	}

	updateLock.Lock()
	defer updateLock.Unlock()

	versionId := domain.RandomId()
	now := domain.TimeNow()
	etag := domain.NewEtag()
	if _, err := createObjectVersionStmt.ExecContext(ctx, versionId, o.ID, cmd.ContentType, size, now, etag); err != nil {
		return nil, fmt.Errorf("unable to create object version: %v", err)
	}
	if _, err := addObjectChunkStmt.ExecContext(ctx, versionId, chunkId, 1); err != nil {
		return nil, fmt.Errorf("unable to persist object chunk record: %v", err)
	}
	if _, err := updateObjectMetadataStmt.ExecContext(ctx, cmd.ContentType, size, etag, versionId, o.ID); err != nil {
		return nil, fmt.Errorf("unable to update object: %v", err)
	}
	if _, err := markObjectVersionDeletedStmt.ExecContext(ctx, o.CurrentVersion); err != nil {
		return nil, fmt.Errorf("unable to set previous object version as deleted")
	}

	triggerPurge()

	return &Object{
		ID:             o.ID,
		Bucket:         o.Bucket,
		Key:            o.Key,
		ContentType:    cmd.ContentType,
		Size:           size,
		Deleted:        o.Deleted,
		ETag:           etag,
		CurrentVersion: versionId,
	}, nil
}

func Write(ctx context.Context, o *Object, w io.Writer) error {
	chunkIds, err := findObjectChunks(ctx, o.CurrentVersion)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Write(chunkId, w); err != nil {
			return err
		}
	}
	return nil
}

func Delete(ctx context.Context, o *Object) error {
	o.Deleted = true
	if _, err := markObjectDeletedStmt.ExecContext(ctx, o.ID); err != nil {
		return fmt.Errorf("unable to update object record: %v", err)
	}
	triggerPurge()
	return nil
}

func findObjectChunks(ctx context.Context, objectId string) ([]string, error) {
	//TODO: when the number of chunks becomes large, this needs to "cursor" its way through
	rows, err := findObjectChunksStmt.QueryContext(ctx, objectId)
	if err != nil {
		return nil, fmt.Errorf("unable to find chunks: %v", err)
	}
	ids := make([]string, 0)
	for rows.Next() {
		var chunkId string
		if err := rows.Scan(&chunkId); err != nil {
			return nil, fmt.Errorf("unable to decode chunk: %v", err)
		}
		ids = append(ids, chunkId)
	}
	return ids, nil
}

func triggerPurge() {
	purgeMutex.Lock()
	purgeFlag = true
	purgeMutex.Unlock()
}

func worker() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		purge()
	}
}

func purge() {
	ctx, cancel := context.WithTimeout(context.Background(), maxPurgeTime)
	purgeContext(ctx)
	defer cancel()
}

func purgeContext(ctx context.Context) {
	if !purgeFlag {
		return
	}
	objectIds, err := getDeletedObjectIds(ctx)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	for i, id := range objectIds {
		select {
		case <-ctx.Done():
			if i > 0 {
				log.Printf("purged %d objects", i)
			}
			return
		default:
			if err := purgeObject(ctx, id); err != nil {
				log.Printf("unable to purge object: %v", err)
				return
			}
		}
	}
	if len(objectIds) > 0 {
		log.Printf("purged %d objects", len(objectIds))
	}

	versionIds, err := getDeletedObjectVersionIds(ctx)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	for i, id := range versionIds {
		select {
		case <-ctx.Done():
			if i > 0 {
				log.Printf("purged %d object versions", i)
			}
		default:
			if err := purgeObjectVersion(ctx, id); err != nil {
				log.Printf("unable to purge object version: %v", err)
				return
			}
		}
	}

	if len(versionIds) > 0 {
		log.Printf("purged %d object versions", len(versionIds))
	}

	purgeMutex.Lock()
	purgeFlag = false
	purgeMutex.Unlock()
}

func getDeletedObjectIds(ctx context.Context) ([]string, error) {
	rows, err := findDeletedObjectsStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to find deleted objects: %v", err)
	}
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("unable to scan object row: %v", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func purgeObject(ctx context.Context, objectId string) error {
	if _, err := markObjectVersionsDeletedStmt.ExecContext(ctx, objectId); err != nil {
		return fmt.Errorf("unable to mark object versions as deleted: %v", err)
	}
	if _, err := deleteObjectStmt.ExecContext(ctx, objectId); err != nil {
		return fmt.Errorf("unable to delete object: %v", err)
	}
	return nil
}

func getDeletedObjectVersionIds(ctx context.Context) ([]string, error) {
	rows, err := findDeletedObjectVersionsStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to find deleted object versions: %v", err)
	}
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("unable to scan object version row: %v", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func purgeObjectVersion(ctx context.Context, versionId string) error {
	chunkIds, err := findObjectChunks(ctx, versionId)
	if err != nil {
		return err
	}
	for _, chunkId := range chunkIds {
		if err := chunk.Delete(ctx, chunkId); err != nil {
			return err
		}
	}
	if _, err := deleteObjectChunksStmt.ExecContext(ctx, versionId); err != nil {
		return fmt.Errorf("unable to delete chunk links: %v", err)
	}

	if _, err := deleteObjectVersionStmt.ExecContext(ctx, versionId); err != nil {
		return fmt.Errorf("unable to delete object version: %v", err)
	}
	return nil
}
