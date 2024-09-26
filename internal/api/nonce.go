// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/nonce"
	"github.com/cfichtmueller/stor/internal/ec"
)

type NonceResponse struct {
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func handleCreateNonce(c jug.Context) {
	o, ok := objectFilter(c)
	if !ok {
		return
	}

	if !c.Request().URL.Query().Has("ttl") {
		handleError(c, ec.InvalidArgument)
		return
	}

	ttl, err := c.IntQuery("ttl")
	if err != nil {
		handleError(c, err)
		return
	}

	n, err := nonce.Create(c, o.Bucket, o.Key, nonce.CreateCommand{
		TTL: time.Duration(ttl) * time.Second,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.RespondCreated(NonceResponse{
		Nonce:     n.ID,
		ExpiresAt: n.ExpiresAt,
	})
}
