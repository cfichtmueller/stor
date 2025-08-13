package object

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
)

func Test_lifecycle(t *testing.T) {
	config.DataDir = os.TempDir()
	ctx := context.Background()
	bucketName := "lifecycle-test"
	db.Configure()
	chunk.Configure()
	Configure()

	initialObjects := countRows(t, objectsTable)
	initialVersions := countRows(t, objectVersionsTable)
	initialObjectChunks := countRows(t, objectChunksTable)
	initialChunks := countRows(t, chunksTable)

	key := uniqueString("o-")

	o, err := Create(ctx, bucketName, CreateCommand{
		Key:         key,
		ContentType: "text/plain",
		Data:        []byte(uniqueString("Hello World ")),
	})
	if err != nil {
		t.Errorf("unable to create object: %v", err)
	}

	expectRows(t, "create object", objectsTable, initialObjects+1)
	expectRows(t, "create object", objectVersionsTable, initialVersions+1)
	expectRows(t, "create object", objectChunksTable, initialObjectChunks+1)
	expectRows(t, "create object", chunksTable, initialChunks+1)

	updated, err := Update(ctx, o, UpdateCommand{
		ContentType: "text/plain",
		Data:        []byte(uniqueString("Pretty new here ")),
	})
	if err != nil {
		t.Errorf("unable to update object: %v", err)
	}

	purge()

	expectRows(t, "update object", objectsTable, initialObjects+1)
	expectRows(t, "update object", objectVersionsTable, initialVersions+1)
	expectRows(t, "update object", objectChunksTable, initialObjectChunks+1)
	expectRows(t, "update object", chunksTable, initialChunks+1)

	if err := Delete(ctx, updated); err != nil {
		t.Errorf("unable to delete object: %v", err)
	}

	purge()

	expectRows(t, "delete object", objectsTable, initialObjects)
	expectRows(t, "delete object", objectVersionsTable, initialVersions)
	expectRows(t, "delete object", objectChunksTable, initialObjectChunks)
	expectRows(t, "delete object", chunksTable, initialChunks)
}

func countRows(t *testing.T, table string) int64 {
	var count int64
	if err := db.QueryRow("SELECT COUNT(*) AS count FROM " + table).Scan(&count); err != nil {
		t.Errorf("unable to count rows in table %s: %v", table, err)
		t.FailNow()
	}
	return count
}

func expectRows(t *testing.T, checkpoint, table string, expected int64) {
	actual := countRows(t, table)
	if actual != expected {
		t.Errorf("Expected %d rows in %s after %s, got %d", expected, table, checkpoint, actual)
		t.FailNow()
	}
}

func uniqueString(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixMilli())
}
