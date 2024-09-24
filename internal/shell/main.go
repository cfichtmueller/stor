// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shell

import (
	"log"

	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/archive"
	"github.com/cfichtmueller/stor/internal/domain/bucket"
	"github.com/cfichtmueller/stor/internal/domain/chunk"
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
}
