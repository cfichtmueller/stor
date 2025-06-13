// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/nonce"
	"github.com/cfichtmueller/stor/internal/ec"
)

type NonceResponse struct {
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func handleCreateNonce(c *srv.Context) *srv.Response {
	o, r := objectFilter(c)
	if r != nil {
		return r
	}

	if !c.HasQuery("ttl") {
		return responseFromError(ec.InvalidArgument)
	}

	ttl, r := c.IntQuery("ttl")
	if r != nil {
		return r
	}

	n, err := nonce.Create(c, o.Bucket, o.Key, nonce.CreateCommand{
		TTL: time.Duration(ttl) * time.Second,
	})
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().Created(NonceResponse{
		Nonce:     n.ID,
		ExpiresAt: n.ExpiresAt,
	})
}
