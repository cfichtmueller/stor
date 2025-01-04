// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/jug"
)

type NodeHandler func(c jug.Context) (e.Node, error)

func renderNode(h NodeHandler) func(c jug.Context) {
	return func(c jug.Context) {
		n, err := h(c)
		if err != nil {
			c.HandleError(err)
			return
		}
		if n == nil {
			return
		}
		if err := n.Render(c.Writer()); err != nil {
			c.HandleError(err)
		}
	}
}
