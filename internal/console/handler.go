// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/srv"
)

type NodeHandler func(c *srv.Context) (e.Node, error)

func renderNode(h NodeHandler) func(c *srv.Context) *srv.Response {
	return func(c *srv.Context) *srv.Response {
		n, err := h(c)
		if err != nil {
			return responseFromError(err)
		}
		return nodeResponse(n)
	}
}
