// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shell

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/archive"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
	"github.com/cfichtmueller/stor/internal/domain/nonce"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/domain/user"
)

func Configure() {
	if len(config.DataDir) == 0 {
		log.Fatal("data dir not set")
	}

	db.Configure()

	user.Configure()
	apikey.Configure()
	session.Configure()
	chunk.Configure()
	bucket.Configure()
	object.Configure()
	archive.Configure()
	nonce.Configure()
}

func Check() bool {
	if len(config.DataDir) == 0 {
		log.Fatal("data dir not set")
	}
	_, err := os.Stat(config.DataDir)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("ERROR: data directory does not exist\n")
		return false
	}
	if !db.Check() {
		return false
	}

	db.Configure()

	if !chunk.Check() {
		return false
	}

	return true
}
